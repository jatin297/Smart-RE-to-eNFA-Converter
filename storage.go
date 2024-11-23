package main

import (
	"database/sql"
	"fmt"
	. "github.com/jatin297/retoenfa/user"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Storage interface {
	CreateUser(user *User) error
	DeleteUser(id int) error
	UpdateUser(user *User) error
	GetUserByID(int) (*User, error)
	GetAllUsers() ([]*User, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStorage, error) {
	connec_string := "user=postgres dbname=postgres password=jatinsalgotra sslmode=disable"
	db, err := sql.Open("postgres", connec_string)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Init() error {
	return s.CreateUserTable()
}

func (s *PostgresStorage) CreateUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users(
		id serial PRIMARY KEY,
		name TEXT,
		email TEXT,
		password TEXT,
		created_at timestamptz NOT NULL DEFAULT NOW()
	)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) CreateUser(user *User) error {

	password, _ := bcrypt.GenerateFromPassword([]byte(user.EncryptedPassword), bcrypt.DefaultCost)
	user.EncryptedPassword = string(password)

	query := `INSERT INTO users(name, email, password, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := s.db.Exec(query, user.Name, user.Email, user.EncryptedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		return err
	}
	d, _ := result.RowsAffected()
	if d == 0 {
		return fmt.Errorf("no rows found with id %d", id)
	}
	return nil
}

func (s *PostgresStorage) UpdateUser(user *User) error {
	return nil
}

func (s *PostgresStorage) GetUserByID(id int) (*User, error) {
	query := `SELECT * FROM users where id = $1`
	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	for row.Next() {
		var user User
		if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.EncryptedPassword, &user.CreatedAt); err != nil {
			return &user, err
		}
		return &user, nil
	}
	return nil, nil
}

func (s *PostgresStorage) GetAllUsers() ([]*User, error) {
	query := `SELECT * FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.EncryptedPassword, &user.CreatedAt); err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}
