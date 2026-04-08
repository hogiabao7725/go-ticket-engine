package hash

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	// Table of test cases
	tests := []struct {
		name     string
		password string
		wantErr  bool
		err      error
	}{
		{
			name:     "Hash Password Successfully",
			password: "1234",
			wantErr:  false,
		},
		{
			name:     "Hash Empty Password",
			password: "",
			wantErr:  true,
			err:      ErrPasswordEmpty,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashed, err := HashPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.err != nil {
					assert.True(t, errors.Is(err, tt.err))
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashed)

				// Check if the hashed password contains the bcrypt prefix
				assert.Contains(t, hashed, "$2")
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	// Pre-hash a password for testing
	validPass := "1234"
	validHashed, _ := HashPassword(validPass)

	tests := []struct {
		name       string
		hashedPass string
		password   string
		wantErr    bool
		err        error
	}{
		{
			name:       "Correct Password",
			hashedPass: validHashed,
			password:   validPass,
			wantErr:    false,
		},
		{
			name:       "Wrong Password",
			hashedPass: validHashed,
			password:   "wrong",
			wantErr:    true,
			err:        ErrInvalidCredentials,
		},
		{
			name:       "Invalid Hashed Password",
			hashedPass: "invalid-hash",
			password:   validPass,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hashedPass, tt.password)
			if tt.wantErr {
				assert.Error(t, err)

				if tt.err != nil {
					assert.True(t, errors.Is(err, tt.err))
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
