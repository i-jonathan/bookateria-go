package account

import (
	"bookateria-api-go/core"
	"bookateria-api-go/log"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

type User struct {
	// For Returning Data, might have to create another struct that is used solely for reading from
	// Seems there's no write only for json or gorm for that matter
	gorm.Model
	UserName  string    `json:"user_name"`
	FirstName string    `json:"first_name" gorm:"not null"`
	LastName  string    `json:"last_name" gorm:"not null"`
	Email     string    `json:"email" gorm:"not null;unique"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	Password  string    `json:"password"`
	LastLogin time.Time `json:"last_login"`
	IsActive  bool      `json:"is_active" gorm:"default:false"`
	IsEmailVerified bool `json:"is_email_verified" gorm:"default:false"`
}

type Profile struct {
	gorm.Model
	//Bio		string	`json:"bio"`
	//Picture	string	`json:"picture"`
	Points int  `json:"points" gorm:"default:20"`
	UserID int  `json:"user_id"`
	User   User `json:"user" gorm:"constraints:OnDelete:CASCADE;not null;unique"`
}

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func GeneratePasswordHash(password string) (string, error) {
	// This function is for generating the hash that would be stored in the database.
	// It takes the password as a string and using argon2id hash algorithm, sends output of
	// the proper format for storing argon2 hashes.

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	config := &PasswordConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
	hash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format for storing argon2id in database: argon2 version, memory, time,
	// number of threads, salt and hash encoded in base 64
	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, config.memory, config.time, config.threads, b64Salt, b64Hash)
	return full, nil
}

func ComparePassword(password, hash string) (bool, error) {
	// This function takes in the password and the hash stored in the database as strings
	// to compare and confirm that the password is correct.
	// Uses constant time compare to prevent timing attacks
	// Returns true, nil if password is correct
	parts := strings.Split(hash, "$")
	config := &PasswordConfig{}

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &config.memory, &config.time, &config.threads)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	config.keyLen = uint32(len(decodedHash))

	comparisonHash := argon2.IDKey([]byte(password), salt, config.time, config.memory, config.threads, config.keyLen)

	return subtle.ConstantTimeCompare(decodedHash, comparisonHash) == 1, nil
}

func commonPasswordValidator(password string) bool {
	// This function makes sure that the password is not in a list of common passwords.
	// The list has been trimmed down to save time. Since no password that is less than 8 characters
	// is to be allowed, the list doesn't include them. Also all numeric passwords are not included.
	// This implements a binary search to check if the password is common.
	// Returns true if the password is common, false if it isn't
	index := sort.Search(len(core.CommonPasswords), func(i int) bool {
		return core.CommonPasswords[i] >= password
	})

	if index <= len(core.CommonPasswords) && core.CommonPasswords[index] == password {
		return true
	}
	return false
}

func SimilarToUser(firstName, lastName, username, password string) bool {
	// Checks if the password is in any case similar to the inputted user information
	// Returns true if password is similar to first name, last name, user name

	containsFirst := strings.Contains(strings.ToLower(password), strings.ToLower(firstName))
	containsLast := strings.Contains(strings.ToLower(password), strings.ToLower(lastName))
	containsUser := strings.Contains(strings.ToLower(password), strings.ToLower(username))

	return containsFirst && containsLast && containsUser
}

func PasswordValidator(password string) bool {
	// Complete password validator. This aggregates all the conditions that a password needs to meet
	// Length, common, uppercase, lowercase and number
	// If the password is good to go, it returns true.
	// And then the password can be hashed then saved.

	var (
		passLen     = len(password) >= 8
		isNotCommon = !commonPasswordValidator(password)
		hasUpper    = false
		hasLower    = false
		hasNumber   = false
	)

	for _, i := range password {
		switch {
		case unicode.IsUpper(i):
			hasUpper = true
		case unicode.IsLower(i):
			hasLower = true
		case unicode.IsNumber(i):
			hasNumber = true
		}
	}

	return passLen && isNotCommon && hasNumber && hasLower && hasUpper
}

func UserDetails(firstName, lastName, email string) (string, string, string, bool) {
	// This attempts to normalize the user details. If they are not empty
	// If empty, returns false
	// Else, returns the details as Title case
	firstName = strings.ReplaceAll(firstName, " ", "")
	lastName = strings.ReplaceAll(lastName, " ", "")
	email = strings.ReplaceAll(email, " ", "")

	if firstName == "" || lastName == "" || email == "" {
		return firstName, lastName, email, false
	}
	firstName = strings.Title(firstName)
	lastName = strings.Title(lastName)
	email = strings.ToLower(email)

	return firstName, lastName, email, true
}

func EmailValidator(email string) bool {
	// This function does (currently) 2 checks on the email to ensure it is correct
	// A regex check and an MX lookup that checks if the domain has MX records
	// The regex check is ridiculously simple because.
	// 1, We are still doing an MX lookup
	// 2, We would still send a verification email. So why make it complex.
	// Returns true if the email is good to go and false otherwise
	re := regexp.MustCompile("^.+@.+\\..+$")
	validity := re.MatchString(email)
	if !validity {
		return false
	}
	parts := strings.Split(email, "@")
	mx, err := net.LookupMX(parts[1])
	log.ErrorHandler(err)
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

func InitDatabase() *gorm.DB {
	viperConfig := core.ReadViper()
	var (
		databaseName = viperConfig.Get("database.name")
		port         = viperConfig.Get("database.port")
		pass         = viperConfig.Get("database.pass")
		user         = viperConfig.Get("database.user")
		host         = viperConfig.Get("database.host")
		ssl          = viperConfig.Get("database.ssl")
	)

	postgresConnection := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s",
		host, port, user, databaseName, pass, ssl)
	db, err := gorm.Open(postgres.Open(postgresConnection), &gorm.Config{})
	log.ErrorHandler(err)

	err = db.AutoMigrate(&User{})
	err = db.AutoMigrate(&Profile{})
	log.ErrorHandler(err)
	return db
}
