package main

import (
	"fmt"
	"time"

	"github.com/pquerna/otp/totp"
)

func main() {
	secret := "V5O7Q7ZX4CPTSWN7R3XGVXLI65ZDDYHY"

	// Generate code for current time
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Generated TOTP Code:", code)
	fmt.Println("(Valid for ~30 seconds, then changes)")
}
