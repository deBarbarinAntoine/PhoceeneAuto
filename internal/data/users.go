package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"PhoceeneAuto/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

var (
	AnonymousUser = &User{}
	UserRole      = struct {
		ADMIN string
		USER  string
	}{
		ADMIN: "ADMIN",
		USER:  "USER",
	}
	UserStatus = struct {
		ACTIVE   string
		INACTIVE string
	}{
		ACTIVE:   "ACTIVE",
		INACTIVE: "INACTIVE",
	}
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   Address   `json:"address"`
	Password  password  `json:"-"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	Shop      string    `json:"shop"`
	Version   int       `json:"-"`
}

type UserSql struct {
	ID        int            `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Name      sql.NullString `json:"name"`
	Email     sql.NullString `json:"email"`
	Phone     sql.NullString `json:"phone"`
	Address   AddressSql     `json:"address"`
	Password  password       `json:"-"`
	Role      sql.NullString `json:"role"`
	Status    sql.NullString `json:"status"`
	Shop      sql.NullString `json:"shop"`
	Version   int            `json:"-"`
}

func (userSql UserSql) toUser() User {
	user := User{
		ID:        userSql.ID,
		CreatedAt: userSql.CreatedAt,
		UpdatedAt: userSql.UpdatedAt,
		Address:   userSql.Address.toAddress(),
		Password:  userSql.Password,
		Version:   userSql.Version,
	}
	if userSql.Name.Valid {
		user.Name = userSql.Name.String
	}
	if userSql.Email.Valid {
		user.Email = userSql.Email.String
	}
	if userSql.Phone.Valid {
		user.Phone = userSql.Phone.String
	}
	if userSql.Role.Valid {
		user.Role = userSql.Role.String
	}
	if userSql.Status.Valid {
		user.Status = userSql.Status.String
	}
	if userSql.Shop.Valid {
		user.Shop = userSql.Shop.String
	}
	return user
}

func EmptyUser() *User {
	return &User{
		Shop:   Shop.HEADQUARTERS,
		Status: UserStatus.ACTIVE,
		Role:   UserRole.USER,
	}
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func IsAdmin(role string) bool {
	return role == UserRole.ADMIN
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainTextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
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
	v.Check(len(password) >= 8, "password", "must not be more than 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "username", "must be provided")
	v.Check(len(user.Name) <= 500, "username", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	user.Address.validate(v)

	v.Check(validator.PermittedValue(user.Role, UserRole.USER, UserRole.USER), "role", fmt.Sprintf("invalid role %s", user.Role))

	v.Check(validator.PermittedValue(user.Status, UserStatus.ACTIVE, UserStatus.INACTIVE), "status", fmt.Sprintf("invalid status %s", user.Status))

	v.Check(validator.PermittedValue(user.Shop, Shop.HEADQUARTERS), "shop", fmt.Sprintf("invalid shop %s", user.Shop))

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	db *sql.DB
}

func (m UserModel) Insert(user *User) error {

	// creating the query
	query := `
		INSERT INTO users (username, email, password_hash, phone, status, shop, street, complement, city, zip_code, state)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, version;`

	// setting the arguments
	args := []any{user.Name, user.Email, user.Password.hash, user.Phone, user.Status, user.Shop}
	args = user.Address.toSQL(args)

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
	err = stmt.QueryRowContext(ctx, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Update(user *User) error {

	// creating the query
	query := `
		UPDATE users
		SET username = $1, email = $2, password_hash = $3, phone = $4, status = $5, shop = $6, street = $9, complement = $10, city = $11, zip_code = $12, state = $13,
		    updated_at = CURRENT_TIMESTAMP, version = version + 1
		WHERE id = $7 AND version = $8
		RETURNING version;`

	// setting the arguments
	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Phone,
		user.Status,
		user.ID,
		user.Version,
	}
	args = user.Address.toSQL(args)

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
	err = stmt.QueryRowContext(ctx, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) Delete(user *User) error {

	// creating the query
	query := `
		DELETE FROM users
		WHERE id = $1 AND version = $2;`

	// setting the arguments
	args := []any{
		user.ID,
		user.Version,
	}

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

	return err
}

func (m UserModel) Exists(id int) (bool, error) {

	// creating the query
	query := `
		SELECT EXISTS (
		SELECT 1 FROM users WHERE id = $1);`

	// setting the timeout for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	// executing the query
	var exists bool
	err = stmt.QueryRowContext(ctx, id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (m UserModel) GetByID(id int) (*User, error) {

	// creating the query
	query := `
		SELECT id, created_at, updated_at, username, email, password_hash, phone, status, shop, street, complement, city, zip_code, state, version
		FROM users
		WHERE id = $1;`

	// setting the userSql variable
	var userSql UserSql

	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// executing the query
	err = stmt.QueryRowContext(ctx, id).Scan(
		&userSql.ID,
		&userSql.CreatedAt,
		&userSql.UpdatedAt,
		&userSql.Name,
		&userSql.Email,
		&userSql.Password.hash,
		&userSql.Phone,
		&userSql.Status,
		&userSql.Shop,
		&userSql.Address.Street,
		&userSql.Address.Complement,
		&userSql.Address.City,
		&userSql.Address.ZIP,
		&userSql.Address.Country,
		&userSql.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	user := userSql.toUser()

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {

	// creating the query
	query := `
		SELECT id, created_at, updated_at, username, email, password_hash, phone, status, shop, street, complement, city, zip_code, state, version
		FROM users
		WHERE email = $1;`

	// setting the userSql variable
	var userSql UserSql

	// setting the timeout context for the query execution
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// preparing the query
	stmt, err := m.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// executing the query
	err = stmt.QueryRowContext(ctx, email).Scan(
		&userSql.ID,
		&userSql.CreatedAt,
		&userSql.UpdatedAt,
		&userSql.Name,
		&userSql.Email,
		&userSql.Password.hash,
		&userSql.Phone,
		&userSql.Status,
		&userSql.Shop,
		&userSql.Address.Street,
		&userSql.Address.Complement,
		&userSql.Address.City,
		&userSql.Address.ZIP,
		&userSql.Address.Country,
		&userSql.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	user := userSql.toUser()

	return &user, nil
}
