package models

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var db *pgx.Conn

func New(dbPool *pgx.Conn) Models {
	db = dbPool

	return Models{
		User:  User{},
		Token: Token{},
	}
}

type Models struct {
	User  User
	Token Token
}

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     Token     `json:"token"`
}

// GetAll gets all users in the database ordered by last_name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT id, email, first_name, last_name, password, created_at, updated_at
		FROM users ORDER BY last_name`

	rows, err := db.Query(ctx, stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetByEmail gets a user from the database that matches the email param
func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT id, email, first_name, last_name, password, created_at, updated_at
		FROM users WHERE email = $1`

	var user User
	row := db.QueryRow(ctx, stmt, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetById gets a user from the database that matches the id param
func (u *User) GetById(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT id, email, first_name, last_name, password, created_at, updated_at
		FROM users WHERE id = $1`

	var user User
	row := db.QueryRow(ctx, stmt, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates the user data in the database
func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		UPDATE users SET
			email = $1,
			first_name = $2,
			last_name = $3,
			updated_at = $4
			WHERE id = $5`

	_, err := db.Exec(ctx, stmt,
		u.Email,
		u.FirstName,
		u.LastName,
		time.Now(),
		u.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete deletes a user from the database
func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		DELETE FROM users
		WHERE id = $1`

	_, err := db.Exec(ctx, stmt, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// Insert inserts a new user in the database
func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, err
	}

	var newID int
	stmt := `
		INSERT INTO users (email, first_name, last_name,
			password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	err = db.QueryRow(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	TokenHash []byte    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Expiry    time.Time `json:"expiry"`
}
