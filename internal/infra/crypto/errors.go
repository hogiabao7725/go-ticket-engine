package crypto

import "errors"

var (
	ErrEmptyInput  = errors.New("crypto: input string is empty")
	ErrMismatched  = errors.New("crypto: hashed value does not match plain text")
	ErrInvalidHash = errors.New("crypto: provided hash is invalid")
)
