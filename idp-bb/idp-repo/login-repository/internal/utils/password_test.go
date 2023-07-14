package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"testing"
)

func TestHashPassword(t *testing.T) {
	type testCase struct {
		name string
		password string
		expectedError bool
	}

	testCases := []testCase{
		{name: "hashing a valid password", password: "testpassword", expectedError: false},
		//{name: "hashing an empty password", password: "", expectedError: true},
		//{name: "hashing a password with only spaces", password: "      ", expectedError: true},
		{name: "hashing a maximum length password", password: string(make([]byte, 72)), expectedError: false},
		//{name: "hashing a password with too many characters", password: string(make([]byte, 73)), expectedError: true},
		{name: "hashing a password with non-ASCII characters", password: "パスワード", expectedError: false},
	}

	log.Println("Testing hashing passwords------------>")

	for i, tc := range testCases {
		log.Printf("Test case #%d: %v\n", i+1, tc.name)

		hashedPassword, err := HashPassword(tc.password)

		if err != nil && !tc.expectedError {
			t.Fatalf("unexpected error hashing password: %v\n", err)
		}

		if err == nil && tc.expectedError {
			t.Fatalf("expected an error for password: '%s', but got nil\n", tc.password)
		}

		if !tc.expectedError {
			err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tc.password))
			if err != nil {
				t.Fatalf("hashed password does not match the original password: %v\n", err)
			}
		}
		log.Printf("Test case #%d is passed\n", i+1)
	}

	log.Println("All test cases are finished.")

}
