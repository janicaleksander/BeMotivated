package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/janicaleksander/BeMotivated/types"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "qwop"

	dbname = "postgres"
)

type Storage interface {
	GetDataBase() (*sql.DB, error)
	CreateAccount(*types.Account) error
	GetAccount(email, pwd string) (*types.Account, error)
	CreateTask(task *types.Task) error
	GetTask(id int) ([]types.Task, error)
	DeleteTask(itemID int) error
	GetTaskByDate(id int, date time.Time) ([]types.Task, error)
	TestChart(id int) error
	IncreasePomodoro(id int) error
	CountDailyStreak(id int) ([]int, error)
	CountCategory(id int) []int
}

type Postgres struct {
	db *sql.DB
}

func (s *Postgres) Init() (error, error, error) {
	return s.createAccountTable(), s.createTaskTable(), s.createPomodoroTable()
}

func NewPostgresDB() (*Postgres, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("dbUser")
	dbPassword := os.Getenv("dbPassword")
	dbHost := os.Getenv("dbHost")
	dbName := os.Getenv("dbName")
	dbPort := os.Getenv("dbPort")

	p, _ := strconv.Atoi(dbPort)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=require",
		dbHost, p, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{db: db}, nil
	/*	err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		databaseUrl := os.Getenv("DATABASE_URL")
		if databaseUrl == "" {
			log.Fatal("DATABASE_URL is not set")
		}
		db, err := sql.Open("postgres", databaseUrl)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v\n", err)
		}
		defer db.Close()

		// Sprawdź połączenie
		err = db.Ping()
		if err != nil {
			log.Fatalf("Cannot ping database: %v\n", err)
		}

		fmt.Println("Successfully connected to the database!")
		return &Postgres{db: db}, nil*/
	/*	serviceURI := os.Getenv("DATABASE_URL")

		conn, _ := url.Parse(serviceURI)
		conn.RawQuery = "sslmode=verify-ca;sslrootcert=ca.pem"

		db, err := sql.Open("postgres", conn.String())

		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		rows, err := db.Query("SELECT version()")
		if err != nil {
			panic(err)
		}

		for rows.Next() {
			var result string
			err = rows.Scan(&result)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Version: %s\n", result)
		}
		return &Postgres{db: db}, nil
	*/
}

func (s *Postgres) createTaskTable() error {
	query := `CREATE TABLE IF NOT EXISTS task (
        user_id  int,
        item_id SERIAL PRIMARY KEY,
        description VARCHAR(255),
    	category VARCHAR(255),
        created_at TIMESTAMP,
        date_day DATE
    )`
	_, err := s.db.Exec(query)
	return err
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

func (s *Postgres) GetDataBase() (*sql.DB, error) {
	if s.db == nil {
		return nil, errors.New("database connection is nil")
	}
	return s.db, nil
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
	if email == "" || pwd == "" {
		return nil, errors.New("error")
	}
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

func (s *Postgres) CreateTask(task *types.Task) error {
	query := `INSERT INTO task (user_id, description, category, created_at,date_day) VALUES ($1, $2, $3,$4,$5) RETURNING item_id`
	var id int
	err := s.db.QueryRow(query, task.UserID, task.Description, task.Category, task.CreatedAt, task.Date).Scan(&id)
	if err != nil {
		return err
	}
	task.ItemID = id
	return nil
}

func (s *Postgres) GetTask(id int) ([]types.Task, error) {
	query := `SELECT * FROM task WHERE user_id=$1`
	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	var slice []types.Task

	for row.Next() {
		task, err := scanIntoTask(row)
		if err != nil {
			return nil, err
		}
		slice = append(slice, *task)
	}
	return slice, nil
}
func (s *Postgres) GetTaskByDate(id int, date time.Time) ([]types.Task, error) {
	query := `SELECT * FROM task WHERE user_id=$1 AND date_day=$2`
	row, err := s.db.Query(query, id, date)
	if err != nil {
		return nil, err
	}
	var slice []types.Task

	for row.Next() {
		task, err := scanIntoTask(row)
		if err != nil {
			return nil, err
		}
		slice = append(slice, *task)
	}
	return slice, nil
}

func (s *Postgres) DeleteTask(itemID int) error {
	query := `DELETE FROM task WHERE item_id=$1`
	_, err := s.db.Query(query, itemID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (s *Postgres) TestChart(id int) error {
	query := `SELECT value FROM temp WHERE id=$1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return err
	}
	var value int
	for rows.Next() {
		rows.Scan(&value)
	}
	fmt.Println(value)
	return nil
}
func (s *Postgres) createPomodoroTable() error {
	query := `CREATE TABLE IF NOT EXISTS pomodoro (
        user_id  int,
        created_at TIMESTAMP
    )`
	_, err := s.db.Exec(query)
	return err
}

func (s *Postgres) IncreasePomodoro(id int) error {

	query := `INSERT INTO pomodoro VALUES ($1,$2)`
	_, err := s.db.Exec(query, id, time.Now())
	if err != nil {
		return err
	}

	return nil
}

func (s *Postgres) CountDailyStreak(id int) ([]int, error) {
	query := `SELECT created_at FROM pomodoro WHERE created_at >= CURRENT_TIMESTAMP - INTERVAL '1 WEEK' AND user_id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	var mon int
	var tue int
	var wen int
	var thu int
	var fri int
	var sat int
	var sun int

	for rows.Next() {
		date, err := scanIntoDate(rows)
		if err != nil {
			return nil, err
		}

		weekday := date.Weekday().String()

		switch weekday {
		case "Monday":
			mon++
		case "Tuesday":
			tue++
		case "Wednesday":
			wen++
		case "Thursday":
			thu++
		case "Friday":
			fri++
		case "Saturday":
			sat++
		case "Sunday":
			sun++
		default:
			fmt.Println("Nieznany dzień:", weekday)
		}
	}

	return []int{mon, tue, wen, thu, fri, sat, sun}, nil

}

func (s *Postgres) CountCategory(id int) []int {
	query := `SELECT category, COUNT(*) AS count FROM task WHERE category IN ('work', 'play', 'training') AND user_id = $1 GROUP BY category;`
	rows, err := s.db.Query(query, id)
	if err != nil {
		fmt.Println(err)
		return nil

	}
	defer rows.Close()
	var play int
	var work int
	var training int

	for rows.Next() {
		var category string
		var count int
		err := rows.Scan(&category, &count)
		fmt.Println(category)
		if err != nil {
			log.Fatal(err)
		}
		if category == "work" {
			work = count
		}
		if category == "play" {
			play = count
		}
		if category == "training" {
			training = count
		}
	}
	return []int{work, play, training}
}

func scanIntoDate(row *sql.Rows) (time.Time, error) {
	var d time.Time
	err := row.Scan(&d)
	return d, err
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
func scanIntoTask(row *sql.Rows) (*types.Task, error) {
	task := new(types.Task)
	err := row.Scan(
		&task.UserID,
		&task.ItemID,
		&task.Description,
		&task.Category,
		&task.CreatedAt,
		&task.Date)
	return task, err
}
