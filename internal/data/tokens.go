package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/deatil/go-encoding/base62"
)

var (
	ErrDuplicateToken = errors.New("duplicate token")
)

// Token represents a token with various fields.
type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// generateToken generates a new token for the given user ID, TTL, and scope.
//
// Parameters:
//
//	userID - The ID of the user associated with the token
//	ttl - The time-to-live duration for the token
//	scope - The scope or type of the token
//
// Returns:
//
//	*Token - A pointer to the newly generated Token instance
//	error - If any error occurs during the process
func generateToken(userID int, ttl time.Duration, scope string) (*Token, error) {

	// creating the token structure with basic data
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// generating the token
	randomBytes := make([]byte, 64)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base62.StdEncoding.EncodeToString(randomBytes)

	// generating the hash to store in the DB
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

// TokenModel represents a model for interacting with tokens in the database.
type TokenModel struct {
	db *sql.DB
}

// New creates a new token and inserts it into the database.
//
// Parameters:
//
//	userID - The ID of the user associated with the token
//	ttl - The time-to-live duration for the token
//	scope - The scope or type of the token
//
// Returns:
//
//	*Token - A pointer to the newly created Token instance
//	error - If any error occurs during the process
func (m TokenModel) New(userID int, ttl time.Duration, scope string) (*Token, error) {

	// generating a new token
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	// saving it in the DB (and regenerate it if it's duplicated)
	err = m.Insert(token)
	if errors.Is(err, ErrDuplicateToken) {
		token, err = generateToken(userID, ttl, scope)
		if err != nil {
			return nil, err
		}

		err = m.Insert(token)
	}

	return token, err
}

// Insert inserts a new token into the database.
//
// Parameters:
//
//	token - The Token instance to insert
//
// Returns:
//
//	error - If any error occurs during the process
func (m TokenModel) Insert(token *Token) error {

	// generating the query
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4);`

	// setting the arguments
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// executing the query
	_, err = stmt.ExecContext(ctx, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "tokens_hash_key"`:
			return ErrDuplicateToken
		default:
			return err
		}
	}

	return nil
}

// DeleteAllForUser deletes all tokens for a given user and scope.
//
// Parameters:
//
//	scope - The scope or type of the token (use "*" to delete all tokens for the user)
//	userID - The ID of the user associated with the tokens
//
// Returns:
//
//	error - If any error occurs during the process
func (m TokenModel) DeleteAllForUser(scope string, userID int) error {

	// generating the query
	query := `
		DELETE FROM tokens
       	WHERE scope = $1 AND user_id = $2;`

	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if scope == "*" {
		// regenerating query
		query = `
			DELETE FROM tokens
       		WHERE user_id = ?;`

		// preparing the query
		stmt, err := m.db.PrepareContext(ctx, query)
		if err != nil {
			return err
		}
		defer stmt.Close()

		// executing the query
		_, err = m.db.ExecContext(ctx, query, userID)
		return err
	}

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// executing the query
	_, err = stmt.ExecContext(ctx, query, scope, userID)
	return err
}

// DeleteExpired deletes all expired tokens from the database.
//
// Returns:
//
//	error - If any error occurs during the process
func (m TokenModel) DeleteExpired() error {

	// generating the query
	query := `
		DELETE FROM tokens
		WHERE expiry < CURRENT_TIMESTAMP;`

	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare query: %w", err)
	}
	defer stmt.Close()

	// executing the query
	_, err = stmt.ExecContext(ctx)
	return err
}
