package hashpassword

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword string, password string) error {
	fmt.Println("hashed password recieved is ", hashedPassword)
	fmt.Println(" password is ", password)
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		fmt.Println("Error recieved by bcrypt comparison", err)
		return err
	}
	return nil
}
