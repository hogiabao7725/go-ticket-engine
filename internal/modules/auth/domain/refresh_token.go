package domain

import "time"

type RefreshToken struct {
	id        string
	userID    string
	tokenHash string
	expiresAt time.Time
	createdAt time.Time
}

func NewRefreshToken(id, userID, tokenHash string, expiresAt, createdAt time.Time) *RefreshToken {
	return &RefreshToken{
		id:        id,
		userID:    userID,
		tokenHash: tokenHash,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}
}

func (r *RefreshToken) ID() string {
	return r.id
}

func (r *RefreshToken) UserID() string {
	return r.userID
}

func (r *RefreshToken) TokenHash() string {
	return r.tokenHash
}

func (r *RefreshToken) ExpiresAt() time.Time {
	return r.expiresAt
}

func (r *RefreshToken) CreatedAt() time.Time {
	return r.createdAt
}
