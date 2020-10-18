package account

import (
	"bookateria-api-go/core"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
	"sort"
	"strings"
	"time"
	"unicode"
)

type User struct {
	// For Returning Data, might have to create another struct that is used solely for reading from
	// Seems there's no write only for json or gorm for that matter
	gorm.Model
	FirstName	string		`json:"first_name"`
	LastName	string		`json:"last_name"`
	Email		string		`json:"email"`
	IsAdmin		bool		`json:"is_admin"`
	Password 	string		`json:"password"`
	LastLogin	time.Time	`json:"last_login"`
}

type Profile struct {
	gorm.Model
	//Bio		string	`json:"bio"`
	//Picture	string	`json:"picture"`
	Points	int		`json:"points" gorm:"default:20"`
	User	User	`json:"user" gorm:"constraints:OnDelete:CASCADE;not null;unique"`
}

type PasswordConfig struct {
	time	uint32
	memory	uint32
	threads	uint8
	keyLen	uint32
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
	// to compare and confirm that the password is correct. Uses constant time compare to prevent timing attacks
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

func PasswordValidator(password string) bool {
	// Complete password validator. This aggregates all the conditions that a password needs to meet
	// Length, common, uppercase, lowercase and number
	// If the password is good to go, it returns true.
	// And then the password can be hashed then saved.
	// TODO include a validator that checks if the password is similar to the user information
	var (
		passLen   = len(password) >= 8
		isNotCommon  = !commonPasswordValidator(password)
		hasUpper  = false
		hasLower  = false
		hasNumber = false
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