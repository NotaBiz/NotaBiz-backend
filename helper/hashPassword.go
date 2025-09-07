// Package helper provides utility functions and structures to assist with
// various operations.
package helper

import "golang.org/x/crypto/bcrypt"

// HashPassword generates a hashed version of the provided plain-text password
// using the bcrypt algorithm.
//
// Parameters:
//   - password: The plain-text password to be hashed.
//
// Returns:
//   - A string containing the hashed password.
//   - An error if the hashing operation fails.
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword compares a hashed password with a plain-text password to check
// if they match. It uses the bcrypt algorithm for comparison.
//
// Parameters:
//   - hashedPassword: The hashed password to compare.
//   - password: The plain-text password to compare against the hashed password.
//
// Returns:
//   - An error if the passwords do not match or if the comparison fails.
//   - Returns nil if the passwords match.
func ComparePassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
