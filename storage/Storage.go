package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/janicaleksander/BeMotivated/types"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwop"
	dbname   = "postgres"
)

type Storage interface {
	GetDataBase() (*sql.DB, error)
	CreateAccount(*types.Account) error
	GetAccount(email, pwd string) (*types.Account, error)
}

type Postgres struct {
	db *sql.DB
}

func (s *Postgres) Init() error {
	return s.createAccountTable()
}

func NewPostgresDB() (*Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
}

func (s *Postgres) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
        id SERIAL PRIMARY KEY,
        nickname VARCHAR(255),
        email VARCHAR(255),
        password VARCHAR(255),
        created_at TIMESTAMP
    )`
	_, err := s.db.Exec(query)
	return err
}

func (s *Postgres) CreateAccount(acc *types.Account) error {
	if err := s.checkUnique("nickname", acc.Nickname); err != nil {
		return err
	}
	if err := s.checkUnique("email", acc.Email); err != nil {
		return err
	}
	query := `INSERT INTO account (nickname, email, password, created_at) VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(query, acc.Nickname, acc.Email, acc.Password, acc.CreatedAt)
	return err

}

func (s *Postgres) GetDataBase() (*sql.DB, error) {
	if s.db == nil {
		return nil, errors.New("database connection is nil")
	}
	return s.db, nil
}

func (s *Postgres) checkUnique(param string, value string) error {
	var query string
	switch param {
	case "nickname":
		query = "SELECT 1 FROM account WHERE nickname = $1"
	case "email":
		query = "SELECT 1 FROM account WHERE email = $1"
	default:
		return errors.New("Unsupported parameter")
	}
	var result int
	err := s.db.QueryRow(query, value).Scan(&result)
	if result == 0 {
		if err == sql.ErrNoRows {
			return nil
		}
		fmt.Println("Error querying database:", err)
		return err
	}

	return errors.New("Not unique")
}
func (s *Postgres) GetAccount(email, pwd string) (*types.Account, error) {
	query := "SELECT * FROM account WHERE email=$1"
	row, err := s.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	account := new(types.Account)
	for row.Next() {
		account, err = scanIntoAccount(row)
	}
	if err != nil {
		log.Default().Print(err)
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(pwd))
	if err != nil {
		return nil, err
	}

	// email and password correct, now generate JWT

	return account, nil

}

func scanIntoAccount(row *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := row.Scan(
		&account.ID,
		&account.Nickname,
		&account.Email,
		&account.Password,
		&account.CreatedAt)

	return account, err
}
