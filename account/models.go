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

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, config.memory, config.time, config.threads, b64Salt, b64Hash)
	return full, nil
}

func ComparePassword(password, hash string) (bool, error) {
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
	index := sort.Search(len(core.CommonPasswords), func(i int) bool {
		return core.CommonPasswords[i] >= password
	})

	if index <= len(core.CommonPasswords) && core.CommonPasswords[index] == password {
		// True means that the password is common
		return true
	}
	return false
}

func PasswordValidator(password string) bool {

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