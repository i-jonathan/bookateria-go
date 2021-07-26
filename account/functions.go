package account

import (
	"bookateriago/core"
	"bookateriago/log"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"math/big"
	"net"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// GenerateOTP Uses crypto/rand package to generate a unique OTP which is used for verification
//  And probably reset password
func generateOTP() string {
	otp, err := rand.Int(rand.Reader, big.NewInt(9999999))
	if err != nil {
		fmt.Println(err)
	}

	return otp.String()
}

// GeneratePasswordHash generates the hash that would be stored in the database.
// It takes the password as a string and using argon2id hashing algorithm, sends output of
// the proper format for storing argon2 hashes.
func generatePasswordHash(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	config := &passwordConfig{
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

// ComparePassword : This function takes in the password and the hash stored in the database as strings
// to compare and confirm that the password is correct.
// Uses constant time compare to prevent timing attacks
// 		return true, nil
// if password is correct
func ComparePassword(password, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	config := &passwordConfig{}

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

// commonPasswordValidator : This function makes sure that the password is not in a list of common passwords.
// The list has been trimmed down to save time. Since no password that is less than 8 characters
// is to be allowed, the list doesn't include them. Also all numeric passwords are not included.
// This implements a binary search to check if the password is common.
// 		return true
// if the password is common, false if it isn't
func commonPasswordValidator(password string) bool {
	index := sort.Search(len(core.CommonPasswords), func(i int) bool {
		return core.CommonPasswords[i] >= password
	})

	if index <= len(core.CommonPasswords) && core.CommonPasswords[index] == password {
		return true
	}
	return false
}

// SimilarToUser checks if the password is in any case similar to the inputted user information
// Returns true if password is similar to first name, last name, user name
func similarToUser(fullName, alias, username, password string) bool {
	containsFull := strings.Contains(strings.ToLower(password), strings.ToLower(fullName))
	containsLast := strings.Contains(strings.ToLower(password), strings.ToLower(alias))
	containsUser := strings.Contains(strings.ToLower(password), strings.ToLower(username))

	return containsFull && containsLast && containsUser
}

// PasswordValidator : Complete password validator. This aggregates all the conditions that a password needs to meet
// Length, common, uppercase, lowercase and number
// If the password is good to go, it returns true.
// And then the password can be hashed then saved.
func passwordValidator(password string) bool {
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

/* userDetails
Checks if string is empty
return true || false
*/
func userDetails(fullName, alias, userName string) bool {
	if strings.Join(strings.Fields(fullName), " ") == "" ||
		strings.Join(strings.Fields(alias), " ") == "" ||
		strings.Join(strings.Fields(userName), " ") == "" {
		return true
	}
	return false
}

/*  emailValidator : This function does (currently) 2 checks on the email to ensure it is correct
A regex check and an MX lookup that checks if the domain has MX records
The regex check is ridiculously simple because.
1, We are still doing an MX lookup
2, We would still send a verification email. So why make it complex.
Returns true if the email is good to go and false otherwise */
func emailValidator(email string) bool {
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

// InitDatabase : Initialize the postgres db and migrate the User and Profile models
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
	err = db.AutoMigrate(&profile{})
	log.ErrorHandler(err)
	return db
}
