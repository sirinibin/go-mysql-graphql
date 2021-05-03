package model

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/jameskeane/bcrypt"
	"gitlab.com/sirinibin/go-mysql-graphql/config"
)

func FindUserByUsername(username string) (user User, err error) {
	err = config.DB.QueryRow("SELECT id,username,password from user where username=?", username).Scan(&user.ID, &user.Username, &user.Password)
	return user, err
}

func FindUserByID(id string) (*User, error) {

	var CreatedAt string
	var UpdatedAt string

	var user User
	err := config.DB.QueryRow("SELECT id,name,username,email,password,created_at,updated_at from user where id=?", id).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &CreatedAt, &UpdatedAt)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05"

	user.CreatedAt, err = time.Parse(layout, CreatedAt)
	if err != nil {
		return nil, err
	}

	user.UpdatedAt, err = time.Parse(layout, UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, err
}

func (user *User) IsEmailExists() (exists bool, err error) {

	var id uint64

	if user.ID != "" {
		//Old Record
		err = config.DB.QueryRow("SELECT id from user where email=? and id!=?", user.Email, user.ID).Scan(&id)
	} else {
		//New Record
		err = config.DB.QueryRow("SELECT id from user where email=?", user.Email).Scan(&id)
	}
	return id != 0, err
}

func (user *User) IsUsernameExists() (exists bool, err error) {

	var id uint64

	if user.ID != "" {
		//Old Record
		err = config.DB.QueryRow("SELECT id from user where username=? and id!=?", user.Username, user.ID).Scan(&id)
	} else {
		//New Record
		err = config.DB.QueryRow("SELECT id from user where username=?", user.Username).Scan(&id)
	}

	return id != 0, err
}

func (user *User) Insert() error {

	res, err := config.DB.Exec("INSERT INTO user(name, username, email, password,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?)", user.Name, user.Username, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s when finding last insert Id", err)
		return err
	}
	user.ID = strconv.FormatInt(id, 10)
	log.Print("user.ID:")
	log.Print(user.ID)
	log.Printf("%d user created ", rows)

	return nil
}

func (user *User) Validate() error {

	emailExists, err := user.IsEmailExists()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if emailExists {
		return errors.New("E-mail is Already in use")
	}

	usernameExists, err := user.IsUsernameExists()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if usernameExists {
		return errors.New("Username is Already in use")
	}

	return nil
}

func HashPassword(password string) string {
	salt, _ := bcrypt.Salt(10)
	hash, _ := bcrypt.Hash(password, salt)
	return hash
}
