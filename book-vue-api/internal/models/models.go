package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"net/http"
	"strings"
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
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Password   string    `json:"password"`
	UserActive int       `json:"user_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Token      Token     `json:"token"`
}

// GetAll gets all users in the database ordered by last_name
func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at,
		CASE WHEN (
			SELECT COUNT(id) FROM tokens AS t
			WHERE user_id = users.id AND t.expiry > NOW()
			) > 0 THEN 1 ELSE 0
		END AS has_token
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
			&user.UserActive,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Token.ID,
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
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM users WHERE email = $1`

	var user User
	row := db.QueryRow(ctx, stmt, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.UserActive,
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
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM users WHERE id = $1`

	var user User
	row := db.QueryRow(ctx, stmt, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.UserActive,
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

// DeleteById deletes a user from the database based in id input
func (u *User) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		DELETE FROM users
		WHERE id = $1`

	_, err := db.Exec(ctx, stmt, id)
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

// ResetPassword resets the users password in the database
func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `UPDATE users SET password = $1 WHERE id = $2`
	_, err = db.Exec(ctx, stmt, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// PasswordMatches checks if the user password matches the password input
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			// invalid password
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
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

// GetByToken gets a token from the database based on token tag
func (t *Token) GetByToken(plainText string) (*Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, user_id, email, token, token_hash, created_at, updated_at, expiry
		FROM tokens WHERE token = $1`

	var token Token
	row := db.QueryRow(ctx, query, plainText)
	err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.Email,
		&token.Token,
		&token.TokenHash,
		&token.CreatedAt,
		&token.UpdatedAt,
		&token.Expiry,
	)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetUserForToken gets a user from the database based on a Token model
func (t *Token) GetUserForToken(token Token) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, email, first_name, last_name, password, created_at, updated_at
		FROM users WHERE id = $1`

	var user User
	row := db.QueryRow(ctx, query, token.UserID)

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

// GenerateToken generates a token and token hash based on user id
func (t *Token) GenerateToken(userID int, ttl time.Duration) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Token = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.Token))
	token.TokenHash = hash[:]

	return token, nil
}

// AuthenticateToken gets a User from the database after authenticating the token
func (t *Token) AuthenticateToken(r *http.Request) (*User, error) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, errors.New("no authorization header received")
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("no valid authorization header received")
	}

	token := headerParts[1]

	if len(token) != 26 {
		return nil, errors.New("token wrong size")
	}

	databaseToken, err := t.GetByToken(token)
	if err != nil {
		return nil, errors.New("no mathing token found")
	}

	if databaseToken.Expiry.Before(time.Now()) {
		return nil, errors.New("expired token")
	}

	user, err := t.GetUserForToken(*databaseToken)
	if err != nil {
		return nil, errors.New("no matching user found")
	}

	return user, nil
}

// Insert inserts a token into the database
func (t *Token) Insert(token Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// delete any existing tokens
	stmt := `
		DELETE FROM tokens
		WHERE user_id = $1`
	_, err := db.Exec(ctx, stmt, token.UserID)
	if err != nil {
		return err
	}

	token.Email = u.Email

	stmt = `
		INSERT INTO tokens (user_id, email, token, token_hash,
			created_at, updated_at, expiry)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = db.Exec(ctx, stmt,
		token.UserID,
		token.Email,
		token.Token,
		token.TokenHash,
		time.Now(),
		time.Now(),
		token.Expiry,
	)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByToken deletes a token from the database with a token tag input
func (t *Token) DeleteByToken(plaintext string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `DELETE FROM tokens WHERE token = $1`
	_, err := db.Exec(ctx, stmt, plaintext)
	if err != nil {
		return err
	}

	return nil
}

// ValidToken checks if the token exists and is valid
func (t *Token) ValidToken(plaintext string) (bool, error) {
	token, err := t.GetByToken(plaintext)
	if err != nil {
		return false, errors.New("no matching token found")
	}

	_, err = t.GetUserForToken(*token)
	if err != nil {
		return false, errors.New("no matching user found")
	}

	if token.Expiry.Before(time.Now()) {
		return false, errors.New("expired token")
	}

	return true, nil
}
