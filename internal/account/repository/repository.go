package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type Repository struct {
	db *sql.DB
}

type Movements struct {
	ID     string
	Date   string
	Amount float64
	Type   string
}

type User struct {
	Alias     string
	FirstName string
	LastName  string
	Email     string
	Password  string
}

var ErrorAlreadyExist = errors.New("already exist")

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// SaveMovements saves the movements in the database in a batch insert
func (r *Repository) SaveMovements(ctx context.Context, movements []Movements, alias string) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("starting transaction: %v", err)
	}
	defer func() {
		if err := recover(); err != nil {
			tx.Rollback() // Rollback the transaction in case of an error
		}
	}()

	// Prepare the insert statement
	insertStmt, err := tx.Prepare("INSERT INTO movements (date, type_movement, amount, alias) VALUES (?, ?, ?, ?);")
	if err != nil {
		return fmt.Errorf("preparing statement: %v", err)
	}
	defer insertStmt.Close()

	// Perform the batch insert
	for _, t := range movements {
		_, err = insertStmt.Exec(t.Date, t.Type, t.Amount, alias)
		if err != nil {
			return fmt.Errorf("inserting transaction: %v", err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}

	return nil
}

// ExistUser checks if the user exists in the database and return a bool value
func (r *Repository) ExistUser(ctx context.Context, alias string) bool {
	var count int
	err := r.db.QueryRow("SELECT count(*) FROM users WHERE alias=?;", alias).Scan(&count)
	if err != nil {
		log.Println("error", err)
		return false

	}

	return count != 0
}

// GetEmail returns the email of the user otherwise an error
func (r *Repository) GetEmail(ctx context.Context, alias string) (string, error) {
	var email string
	err := r.db.QueryRow("SELECT email FROM users WHERE alias=?;", alias).Scan(&email)
	if err != nil {
		return "", err
	}
	return email, nil
}

// Save inserts a new movement in the user account
func (r *Repository) SaveUser(ctx context.Context, u User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users(alias,first_name,last_name,email,password) VALUES(?,?,?,?,?);",
		u.Alias, u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		if v, ok := err.(*mysql.MySQLError); ok {
			fmt.Println("error save user", v.Number)
			if v.Number == 1062 {
				return ErrorAlreadyExist
			}
		}

		return err
	}

	return nil
}
