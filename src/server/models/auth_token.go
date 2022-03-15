package models

import (
	"database/sql"
	"time"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type AuthToken struct {
	Token   string
	UserId  string
	Expires time.Time
}

func CreateAuthToken(tx *db.Tx, userId string) (*AuthToken, error) {
	token, err := id.GenerateId(id.IdTypeToken)
	if err != nil {
		return nil, err
	}

	authToken := AuthToken{Token: token, UserId: userId}

	statement := `INSERT INTO auth_tokens (token, user_id) VALUES ($1, $2) RETURNING expires`
	err = tx.Tx.QueryRow(statement, token, userId).Scan(&authToken.Expires)
	if err != nil {
		return nil, err
	}

	return &authToken, err
}

func FetchAuthToken(tx *db.Tx, token string) (*AuthToken, error) {
	authToken := AuthToken{Token: token}
	statement := `SELECT user_id, expires FROM auth_tokens WHERE token = $1`
	err := tx.Tx.QueryRow(statement, token).Scan(&authToken.UserId, &authToken.Expires)
	if err != nil {
		return nil, err
	}
	if authToken.Expires.Before(time.Now()) {
		statement = `DELETE FROM auth_tokens WHERE token = $1`
		_, err = tx.Tx.Exec(statement, token)
		if err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	return &authToken, nil
}

func (authToken AuthToken) Extend(tx *db.Tx) error {
	statement := `UPDATE auth_tokens SET expires = $1 WHERE token = $2`
	_, err := tx.Tx.Exec(statement, time.Now().AddDate(0, 0, 7), authToken.Token)
	return err
}

func (authToken AuthToken) Delete(tx *db.Tx) error {
	statement := `DELETE FROM auth_tokens WHERE token = $1`
	_, err := tx.Tx.Exec(statement, authToken.Token)
	return err
}

func DeleteAllAuthTokensByUserId(tx *db.Tx, userId string) error {
	statement := `DELETE FROM auth_tokens WHERE user_id = $1`
	_, err := tx.Tx.Exec(statement, userId)
	return err
}
