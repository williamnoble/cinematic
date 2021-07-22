package data

import (
	"database/sql"
	"time"
)

// Models acts as a container for different data models (Movies & Users)
type Models struct {
	Movies MovieModel
	Users  UserModel
	Tokens TokenModel
}

// NewModels returns an instance of Models which holds all our data models.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
	}
}

// Movie describes an individual film entry within Movies database
type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`
	Runtime   Runtime   `json:"runtime,omitempty"` // declare in runtime.go
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type User struct {
	ID        int64     `json:"id""`
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
