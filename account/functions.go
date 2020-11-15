package account

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func GenerateOTP() string {
	otp, err := rand.Int(rand.Reader, big.NewInt(9999999))
	if err != nil {
		fmt.Println(err)
	}

	return otp.String()
}
