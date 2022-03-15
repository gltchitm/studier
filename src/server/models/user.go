package models

import (
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type User struct {
	UserId   string
	Email    string
	Username string
	Password string
	Verified bool
}

func CreateUser(tx *db.Tx, email string, username string, password string) (*User, error) {
	userId, err := id.GenerateId(id.IdTypeRegular)
	if err != nil {
		return nil, err
	}

	user := User{UserId: userId, Email: email, Username: username, Password: password}

	statement := `INSERT INTO users (user_id, email, username, password) VALUES ($1, $2, $3, $4)`
	_, err = tx.Tx.Exec(statement, userId, email, username, password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (user *User) SetVerified(tx *db.Tx, verified bool) error {
	statement := `UPDATE users SET verified = $1 WHERE user_id = $2`
	_, err := tx.Tx.Exec(statement, verified, user.UserId)
	user.Verified = true
	return err
}

func (user *User) SetPassword(tx *db.Tx, password string) error {
	statement := `UPDATE users SET password = $1 WHERE user_id = $2`
	_, err := tx.Tx.Exec(statement, password, user.UserId)
	user.Password = password
	return err
}

func (user *User) SetEmail(tx *db.Tx, email string) error {
	statement := `UPDATE users SET email = $1 WHERE user_id = $2`
	_, err := tx.Tx.Exec(statement, email, user.UserId)
	user.Email = email
	return err
}

func FetchUserByEmail(tx *db.Tx, email string) (*User, error) {
	user := User{Email: email}
	statement := `SELECT user_id, username, password, verified FROM users WHERE email = $1`
	row := tx.Tx.QueryRow(statement, email)
	err := row.Scan(&user.UserId, &user.Username, &user.Password, &user.Verified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUserByUsername(tx *db.Tx, username string) (*User, error) {
	user := User{Username: username}
	statement := `SELECT user_id, email, password, verified FROM users WHERE username = $1`
	row := tx.Tx.QueryRow(statement, username)
	err := row.Scan(&user.UserId, &user.Email, &user.Password, &user.Verified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FetchUserById(tx *db.Tx, userId string) (*User, error) {
	user := User{UserId: userId}
	statement := `SELECT email, username, password, verified FROM users WHERE user_id = $1`
	err := tx.Tx.QueryRow(statement, userId).Scan(&user.Email, &user.Username, &user.Password,
		&user.Verified)
	return &user, err
}

func (user *User) Delete(tx *db.Tx) error {
	statement := `DELETE FROM users WHERE user_id = $1`
	_, err := tx.Tx.Exec(statement, user.UserId)
	return err
}
