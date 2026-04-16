package identifier

import "github.com/google/uuid"

type UUID struct{}

func (UUID) Generate() string {
	return uuid.New().String()
}
