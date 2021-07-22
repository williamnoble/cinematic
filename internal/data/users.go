package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"greenlight/internal/validator"
	"time"
)

var (
	ErrDuplicatedEmail = errors.New("error duplcicated")
)

func (m *UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	// We write the user.Password.hash, ignoring the user.Password.plaintext
	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// we mutate the original struct(User)
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicatedEmail
		default:
			return err
		}
	}

	return nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated, version 
	FROM users
	WHERE email=$1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash, // Must return hash not password after refactor
		&user.Activated,
		&user.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (m *UserModel) Update(user *User) error {
	query := `
	UPDATE users SET 
	name = $1, email = $2, password_hash = $3, activated = $4, 
	version = version + 1
	WHERE id = $5 and version = $6
	RETURNING version`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicatedEmail
		case errors.Is(err, sql.ErrNoRows): // returne by .Scan
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

type password struct {
	plaintext *string // compare nil vs "" for plaintext
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "password must be longer than 8 characters")
	v.Check(len(password) <= 72, "password", "password must be shorter than 18 characters")
}

//ValidateUser takes a new Validator and the User, checking that input is valid and returning a map of errors if a
// problem is found.
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "user", "A user name must be provided")
	v.Check(len(user.Name) <= 500, "user", "must not be longer than 500 bytes long")
	ValidateEmail(v, user.Email)
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil {
		panic("missing hash for user")
	}
}

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// GetForToken recieved the plaintext input token from user. Obtain the SHA256 Hash of this token for
	// compaison to the one contained in the user table.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	// Inner Join: Join both tables where the user is found in the token table. Scope is Authorization and
	// check that we have not exceeded the TTL of the token. Todo: Delete old tokens!

	query := `
	SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
	FROM users
	INNER JOIN tokens
	ON users.ID = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3
	`

	// Again, tokenhash returns a [32]byte, convert to a slice => [32]byte =>[:]
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
