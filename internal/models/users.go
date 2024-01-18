package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	verificationHash, err := bcrypt.GenerateFromPassword([]byte(email+password), 12)
	if err != nil {
		return err
	}

	statement := `INSERT INTO users (email, hashed_password, created, active, verification_hash)
    VALUES(?, ?, UTC_TIMESTAMP(), false, ?)`

	_, err = m.DB.Exec(statement, email, string(hashedPassword), string(verificationHash))
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"

	err := m.DB.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}

func (m *UserModel) Active(id int) (bool, error) {
	var active bool

	stmt := "SELECT active FROM users WHERE id = ?"

	err := m.DB.QueryRow(stmt, id).Scan(&active)
	return active, err
}

func (m *UserModel) GetByVerificationHash(hash string) (int, error) {
	var id int

	stmt := "SELECT id FROM users WHERE verification_hash = ?"

	err := m.DB.QueryRow(stmt, hash).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNoRecord
		} else {
			return 0, err
		}
	}
	return id, err
}

func (m *UserModel) GetVerificationHashByEmail(email string) (string, error) {
	var hash string

	stmt := "SELECT verification_hash FROM users WHERE email = ?"

	err := m.DB.QueryRow(stmt, email).Scan(&hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return hash, ErrNoRecord
		} else {
			return hash, err
		}
	}
	return hash, err
}

func (m *UserModel) Activate(id int) error {
	statement := `UPDATE users SET active = 1 where id = ?`

	_, err := m.DB.Exec(statement, id)

	return err
}
