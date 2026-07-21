package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// bcrypt.GenerateFromPassword: generates a hashed password from the given password and cost
	// The cost parameter determines the computational complexity of the hashing algorithm
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "There was an error hashing the password", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(password, hashedPassword string) bool {
	// bcrypt.CompareHashAndPassword: compares a hashed password with a plaintext password
	// It returns nil if the passwords match, or an error if they don't
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
