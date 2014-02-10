package dbconnection

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"os"
	"time"
	"webserver/models"
)

const (
	db_name           string = "webserver"
	connection_format string = "user=%s dbname=%s sslmode=disable password=%s host=%s"
	get_user          string = "SELECT _id FROM users WHERE username=$1 AND password_hash=$2"
	insert_user       string = "INSERT INTO users (username, password_hash, created_at) VALUES ($1,$2,$3) RETURNING _id"
	insert_user_file  string = "INSERT INTO users_files (title, path, user_id, size, category) VALUES ($1,$2,$3,$4,$5)"
	get_user_file     string = "SELECT title, path FROM users_files WHERE category = $1 LIMIT $2 OFFSET $3"
)

func stablishConnection() (*sql.DB, error) {
	user := os.Getenv("POSTGRESQL_USER")
	pass := os.Getenv("POSTGRESQL_PASS")
	host := os.Getenv("PGHOST")
	connection_params := fmt.Sprintf(connection_format, user, db_name, pass, host)
	db, err := sql.Open("postgres", connection_params)

	if err != nil {
		println("error open")
		return nil, err
	}
	return db, nil
}

func encryptPassword(password string) string {
	hash := md5.New()
	io.WriteString(hash, password)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func VerifyUser(user *models.User) (int64, error) {
	db, err := stablishConnection()
	if err != nil {
		return 0, err
	}
	defer db.Close()
	var user_id int64
	err = db.QueryRow(get_user, user.Username,
		encryptPassword(*user.Password)).Scan(&user_id)
	switch {
	case err == sql.ErrNoRows:
		return 0, errors.New("No user with that ID.")
	case err != nil:
		return 0, err
	}
	return user_id, nil
}

func InsertUser(user *models.User, callback func(int64)) error {
	db, err := stablishConnection()
	if err != nil {
		println("stablishConnection")
		return err
	}
	defer db.Close()

	var transaction *sql.Tx
	transaction, err = db.Begin()
	if err != nil {
		println("transaction")
		return err
	}

	var id int64
	err = transaction.QueryRow(insert_user, user.Username,
		encryptPassword(*user.Password), time.Now()).Scan(&id)
	if err != nil {
		println("transaction.Exec")
		return err
	}

	transaction.Commit()

	if callback != nil {
		callback(id)
	}

	return nil
}

func InsertUserFile(file *models.UserFile) {
	db, err := stablishConnection()
	if err != nil {
		panic("stablishConnection")
		return
	}
	defer db.Close()
	var transaction *sql.Tx
	transaction, err = db.Begin()
	if err != nil {
		panic("transaction")
		return
	}

	_, err = transaction.Exec(insert_user_file,
		file.Title,
		file.Path,
		file.UserId,
		file.Size,
		file.Category)
	if err != nil {
		panic("transaction.Exec " + err.Error())
		return
	}
	transaction.Commit()

}

func GetUsersFiles(limit, offset int, category string) ([]*models.UserFile, error) {
	db, err := stablishConnection()
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	if category == "" {
		category = "animals"
	}

	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query(get_user_file, category, limit, offset)

	if err != nil {
		return nil, err
	}

	f_arr := make([]*models.UserFile, 0, limit)
	var path, title string
	var u_file *models.UserFile
	for rows.Next() {
		if err = rows.Scan(&title, &path); err != nil {
			panic(err)
		}
		if path != "" && title != "" {
			u_file = &models.UserFile{Path: path, Title: title}
			f_arr = append(f_arr, u_file)
		}
	}
	return f_arr, nil
}
