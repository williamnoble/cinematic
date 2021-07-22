package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"greenlight/internal/validator"
	"time"
)

// it is ESSENTIAL we reference crypto/rand and not math/rand.

const ScopeActivation = "activation"
const ScopeAuthentication = "authentication"

// Token allows a user to both Activate and Authenticate depending on scope.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"` // references User.ID on Users table
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlainText string) {
	v.Check(tokenPlainText != "", "token", "must be provided")
	v.Check(len(tokenPlainText) == 26, "token", "must be 26 bytes long")
}

// Make a random Key, Convert to Base32, Generate 256 Hash for Token table.
// Scope: Activation OR Authentication. TTL 3 days for Activation, Auth = 30mins.
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	plainText := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// sha256 returns a BYTE Array hash with a fixed 32 byte size. Converted to slice below.
	hash := sha256.Sum256([]byte(plainText))

	token := &Token{
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
		Plaintext: plainText,
		Hash:      hash[:], // Convert byte array ([32]byte) to slice
	}

	return token, nil
}

//New is a function which generates a new token and then, assuming no error, inserts it to the db.
func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

//Insert adds a token to the tokens table, it stores a SHA256 Hash of the plaintext token
// and a scope indicating whether we are authorizing or authenticating a user.
func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute query without returning any rows
	_, err := m.DB.ExecContext(ctx, query, args...) // ignore Result(lastInsertID)
	return err
}

// DeleteAllForUser Delete all tokens for User given their User.ID and a scope (which may be used
// to indicate authorization or authentication of users).
func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}
