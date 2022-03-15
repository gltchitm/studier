package models

import (
	"database/sql"
	"time"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/email"
	"github.com/gltchtim/studier/server/id"
)

type forgotPasswordToken struct {
	Token   string
	UserId  string
	Expires time.Time
}

const forgotPasswordTokenLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const forgotPasswordTokenLength = 24

func CreateForgotPasswordToken(tx *db.Tx, userId string) (*forgotPasswordToken, error) {
	token, err := id.GenerateIdCustom(forgotPasswordTokenLetters, forgotPasswordTokenLength)
	if err != nil {
		return nil, err
	}

	forgotPasswordtoken := forgotPasswordToken{Token: token, UserId: userId}

	statement := `INSERT INTO forgot_password_tokens (token, user_id) VALUES ($1, $2) RETURNING expires`
	err = tx.Tx.QueryRow(statement, token, userId).Scan(&forgotPasswordtoken.Expires)
	if err != nil {
		return nil, err
	}

	return &forgotPasswordtoken, err
}

func FetchForgotPasswordToken(tx *db.Tx, token string) (*forgotPasswordToken, error) {
	forgotPasswordToken := forgotPasswordToken{Token: token}
	statement := `SELECT user_id, expires FROM forgot_password_tokens WHERE token = $1`
	err := tx.Tx.QueryRow(statement, token).Scan(
		&forgotPasswordToken.UserId,
		&forgotPasswordToken.Expires,
	)
	if err != nil {
		return nil, err
	}
	if forgotPasswordToken.Expires.Before(time.Now()) {
		statement = `DELETE FROM forgot_password_tokens WHERE token = $1`
		_, err = tx.Tx.Exec(statement, forgotPasswordToken.Token)
		if err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	return &forgotPasswordToken, nil
}

func (forgotPasswordToken *forgotPasswordToken) SendEmail(tx *db.Tx) error {
	user, err := FetchUserById(tx, forgotPasswordToken.UserId)
	if err != nil {
		return err
	}

	err = email.SendEmail(user.Email, "Reset Studier Password", "Hello "+user.Username+",\n\n"+
		"Use this code to reset your Studier password:\n"+forgotPasswordToken.Token+"\n\n"+
		"If you did not request to reset your Studier password, you can ignore this email.\n\n"+
		"Thank you.")

	return err
}

func (forgotPasswordToken *forgotPasswordToken) Delete(tx *db.Tx) error {
	statement := `DELETE FROM forgot_password_tokens WHERE token = $1`
	_, err := tx.Tx.Exec(statement, forgotPasswordToken.Token)
	return err
}
