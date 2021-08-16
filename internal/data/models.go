package data

import (
	"database/sql"
	"time"
)

// Models acts as a container to wrap distinct models encapsulating the model definitions.
type Models struct {
	Movies      MovieModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

// NewModels returns an instance of Models which holds all our data models.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Users:       UserModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Permissions: PermissionModel{DB: db},
	}
}

// Token allows a user to both Activate and Authenticate depending on scope.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"` // references User.ID on Users table
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// Movie describes an individual film entry within the movies table.
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` // declare in runtime.go
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

// User describes a single user within the users table.
type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"version"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// UserModel holds the database pool for the users table.
type UserModel struct {
	DB *sql.DB
}

// MovieModel holds the datbase pool for the movies table.
type MovieModel struct {
	DB *sql.DB
}

// PermissionModel hols the database pool which is responsible for managing permissions for a given entity.
type PermissionModel struct {
	DB *sql.DB
}

type TokenModel struct {
	DB *sql.DB
}
