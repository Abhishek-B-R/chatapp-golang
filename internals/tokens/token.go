package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeAuth = "authentication"
)

type Token struct{
	ID int64 `json:"-"`
	UserID int64 `json:"-"`
	Plaintext string `json:"token"`
	Hash []byte `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func GenerateToken(userID int64, ttl time.Duration) (*Token, error){
	token := &Token{
		UserID: userID,
		ExpiresAt: time.Now().Add(ttl),
	}

	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]
	return token, nil
}