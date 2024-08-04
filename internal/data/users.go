package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"ghostprotocols.pk/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	EmailVerfied  bool      `json:"email_verified"`
	Phone         string    `json:"phone"`
	PhoneVerified string    `json:"phone_verified"`
	Password      password  `json:"-"`
	ProfilePic    string    `json:"profile_pic"`
	City          int64     `json:"city"`
	DateJoined    time.Time `json:"date_joined"`
	ListingLimit  int64     `json:"listing_limit"`
	FeaturedLimit int64     `json:"featured_limit"`
	Version       int64     `json:"-"`
	IsDealer      bool      `json:"is_dealer"`
}

type Dealer struct {
	UserID  int64  `json:"user_id,omitempty"`
	Address string `json:"address,omitempty"`
	Timings string `json:"timings,omitempty"`
	Version int64  `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

// Errors Related to User Models
var (
	ErrDuplicateEmail = errors.New("this email already exists")
	ErrDuplicatePhone = errors.New("this phone already exists")
	ErrDuplicateUser  = errors.New("dealer already exists")
)

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
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

// VALIDATORS
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePhone(v *validator.Validator, phone string) {
	v.Check(phone != "", "phone", "must be provided")
	v.Check(validator.Matches(phone, validator.PhoneRX), "phone", "must be a valid phone number")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(len(user.Phone) == 10, "phone", "must be valid phone")

	ValidateEmail(v, user.Email)
	ValidatePhone(v, user.Phone)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func ValidateDealer(v *validator.Validator, user *User, dealer *Dealer) {
	ValidateUser(v, user)
	v.Check(user.City != 0, "city", "city must be provided")
	v.Check(dealer.Address != "", "address", "address must be provided")
	v.Check(dealer.Timings != "", "timings", "timings must be provided")

}

// SQL STATEMENTS

func (m UserModel) InsertUser(user *User) error {

	query := `
	INSERT INTO users (name, email, password_hash, phone, city)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id, date_joined, version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Phone,
		sql.NullInt64{Int64: user.City, Valid: user.City != 0},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.DateJoined, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_phone_key"`:
			return ErrDuplicatePhone
		default:
			return err

		}
	}

	return nil
}

func (m UserModel) InsertDealer(user *User, dealer *Dealer) error {
	m.InsertUser(user)
	dealer.UserID = user.ID

	query := `
	INSERT INTO dealers (user_id, address, timings)
	VALUES ($1, $2, $3)
	RETURNING address, timings;
	`

	args := []any{dealer.UserID, dealer.Address, dealer.Timings}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&dealer.Address, &dealer.Timings)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "dealers_user_id_key"`:
			return ErrDuplicateEmail
		default:
			return err

		}
	}

	return nil
}

func (m UserModel) GetUser(id int64) (*User, error) {
	query := `
	SELECT id, date_joined, name, email, email_verified, 
	phone, phone_verified, password_hash, listing_limit, 
	featured_limit, profile_pic, city,version
	FROM users
	WHERE id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.DateJoined,
		&user.Name,
		&user.Email,
		&user.EmailVerfied,
		&user.Phone,
		&user.PhoneVerified,
		&user.Password.hash,
		&user.ListingLimit,
		&user.FeaturedLimit,
		&user.ProfilePic,
		&user.City,
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

func (m UserModel) GetDealer(id int64) (*Dealer, error) {
	query := `
	SELECT user_id, address, timings, version
	FROM dealers
	WHERE user_id = $1`

	var dealer Dealer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&dealer.UserID,
		&dealer.Address,
		&dealer.Timings,
		&dealer.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &dealer, nil
}

func (m UserModel) UpdateProfilePic(id int64, profile string) error {
	query := `
	UPDATE users
	SET profile_pic = $1
	WHERE users.id = $2
	`

	args := []any{
		profile,
		id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m UserModel) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, email_verified = $3, 
	phone = $4, phone_verified = $5, password_hash = $6, profile_pic = $7,
	city = $10,
	version = version + 1
	WHERE id = $8 AND version = $9
	RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.EmailVerfied,
		user.Phone,
		user.PhoneVerified,
		user.Password.hash,
		sql.NullString{String: user.ProfilePic, Valid: user.ProfilePic != "0"},
		user.ID,
		user.Version,
		user.City,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_phone_key"`:
			return ErrDuplicatePhone
		default:
			return err

		}
	}

	return nil
}

func (m UserModel) UpdateDealer(user *User, dealer *Dealer) error {

	dealer.UserID = user.ID

	query := `
	UPDATE dealers
	SET user_id = $1, address = $2, timings = $3, version = version + 1
	WHERE id = $1 AND version = $4
	RETURNING version;
	`

	args := []any{dealer.UserID, dealer.Address, dealer.Timings, dealer.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "dealers_user_id_key"`:
			return ErrDuplicateEmail
		default:
			return err

		}
	}

	return nil
}

func (m UserModel) GetByPhone(phone string) (*User, error) {
	query := `
	SELECT id, date_joined, name, email, email_verified, 
	phone, phone_verified, password_hash, listing_limit, 
	featured_limit, version
	FROM users
	WHERE phone = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, phone).Scan(
		&user.ID,
		&user.DateJoined,
		&user.Name,
		&user.Email,
		&user.EmailVerfied,
		&user.Phone,
		&user.PhoneVerified,
		&user.Password.hash,
		&user.ListingLimit,
		&user.FeaturedLimit,
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

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
	SELECT id, date_joined, name, email, email_verified, 
	phone, phone_verified, password_hash, listing_limit, 
	featured_limit, version
	FROM users
	WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.DateJoined,
		&user.Name,
		&user.Email,
		&user.EmailVerfied,
		&user.Phone,
		&user.PhoneVerified,
		&user.Password.hash,
		&user.ListingLimit,
		&user.FeaturedLimit,
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

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
	SELECT id, date_joined, name, email, email_verified, 
	phone, phone_verified, password_hash, listing_limit, 
	featured_limit, version FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.DateJoined,
		&user.Name,
		&user.Email,
		&user.EmailVerfied,
		&user.Phone,
		&user.PhoneVerified,
		&user.Password.hash,
		&user.ListingLimit,
		&user.FeaturedLimit,
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
