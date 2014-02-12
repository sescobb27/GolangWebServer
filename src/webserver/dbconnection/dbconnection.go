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
	get_user_file     string = `SELECT ut.title, ut.path FROM
		users_files as ut join (
			file_tags as ft join tags as t
				on (ft.tag_id = t._id and t.name = $1)
			) as jft
		on (jft.file_id = ut._id) LIMIT $2 OFFSET $3`
	insert_file      string = "insert into users_files (title, path, user_id, size) values ($1,$2,$3,$4) returning _id"
	insert_tag       string = "insert into tags (name) values ($1) returning _id"
	insert_file_tags string = "insert into file_tags (file_id, tag_id) values ($1,$2)"
	get_all_tags     string = "select name from tags"
	get_tags         string = "select * from tags"
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

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
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
	assertNoError(err)
	defer db.Close()

	var transaction *sql.Tx
	transaction, err = db.Begin()
	assertNoError(err)

	var id int64
	err = transaction.QueryRow(insert_user, user.Username,
		encryptPassword(*user.Password), time.Now()).Scan(&id)
	assertNoError(err)

	transaction.Commit()

	if callback != nil {
		callback(id)
	}

	return nil
}

func InsertUserFile(file *models.UserFile) {
	db, err := stablishConnection()
	assertNoError(err)

	getTags := make(chan map[string]int64)
	var tags map[string]int64
	go func() {
		tags = getTagsAsObjects()
		// send tags via channel for continue
		getTags <- tags
	}()

	defer db.Close()
	var transaction *sql.Tx
	transaction, err = db.Begin()
	assertNoError(err)

	var fid, tid int64
	err = transaction.QueryRow(insert_file,
		file.Title,
		file.Path,
		file.UserId,
		file.Size).Scan(&fid)

	assertNoError(err)
	transaction.Commit()

	// wait for all the tags
	<-getTags
	fmt.Println(tags)

label:
	for _, tag := range file.Tags {
		transaction, err = db.Begin()
		for k, v := range tags {
			if tag == k {
				_, err = transaction.Exec(insert_file_tags, fid, v)
				assertNoError(err)
				transaction.Commit()
				continue label
			}
		}
		err = transaction.QueryRow(insert_tag, tag).Scan(&tid)
		assertNoError(err)
		_, err = transaction.Exec(insert_file_tags, fid, tid)
		assertNoError(err)
		transaction.Commit()
	}
}

func GetUsersFiles(limit, offset int, tag string) ([]*models.UserFile, error) {
	db, err := stablishConnection()
	assertNoError(err)

	if limit <= 0 {
		limit = 10
	}

	if offset < 0 {
		offset = 0
	}

	if tag == "" {
		tag = "animals"
	}

	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query(get_user_file, tag, limit, offset)

	assertNoError(err)

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

func GetTags() []string {
	db, err := stablishConnection()
	assertNoError(err)
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query(get_all_tags)
	assertNoError(err)

	tags := make([]string, 0)
	var name string
	for rows.Next() {
		if err = rows.Scan(&name); err != nil {
			panic(err)
		}
		tags = append(tags, name)
	}
	return tags
}

func getTagsAsObjects() map[string]int64 {
	db, err := stablishConnection()
	assertNoError(err)
	defer db.Close()
	var rows *sql.Rows
	rows, err = db.Query(get_tags)
	assertNoError(err)

	tags := make(map[string]int64)
	var name string
	var id int64
	for rows.Next() {
		if err = rows.Scan(&id, &name); err != nil {
			panic(err)
		}
		tags[name] = id
	}
	return tags
}
