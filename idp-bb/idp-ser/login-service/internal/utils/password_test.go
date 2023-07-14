package utils

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockRequest struct {
	password string
	hash	 string
}

func TestCheckPassword (t *testing.T) {
	pass := "mock123"
	mockHash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	cases := []struct {
		name	string
		want	error
		request *mockRequest
	}{
		{
			name: "Password and hash valid",
			want: nil,
			request: &mockRequest{
				password: pass,
				hash: string(mockHash),
			},
		},

		{
			name: "Password and hash invalid",
			want: bcrypt.ErrMismatchedHashAndPassword,
			request: &mockRequest{
				password: "wrong123",
				hash: string(mockHash),
			},
		},
	}

	for _, tt := range cases {
		doerr := bcrypt.CompareHashAndPassword([]byte(tt.request.hash), []byte(tt.request.password))
		if doerr != tt.want {
			t.Fatalf("expected return %t, got %t", tt.want, doerr)
		}
	}
}