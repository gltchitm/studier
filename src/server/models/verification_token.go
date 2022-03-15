package models

import (
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/email"
	"github.com/gltchtim/studier/server/id"
)

type VerificationToken struct {
	Token  string
	UserId string
}

const verificationTokenLetters = "ABCDEF0123456789"
const verificationTokenLength = 8

func CreateVerificationToken(tx *db.Tx, userId string) (*VerificationToken, error) {
	token, err := id.GenerateIdCustom(verificationTokenLetters, verificationTokenLength)
	if err != nil {
		return nil, err
	}

	verificationToken := VerificationToken{Token: token, UserId: userId}

	statement := `INSERT INTO verification_tokens (token, user_id) VALUES ($1, $2)`
	_, err = tx.Tx.Exec(statement, token, userId)
	if err != nil {
		return nil, err
	}
	return &verificationToken, err
}

func FetchVerificationToken(tx *db.Tx, userId string) (*VerificationToken, error) {
	verificationToken := VerificationToken{UserId: userId}
	statement := `SELECT token FROM verification_tokens WHERE user_id = $1`
	err := tx.Tx.QueryRow(statement, userId).Scan(&verificationToken.Token)
	if err != nil {
		return nil, err
	}
	return &verificationToken, nil
}

func (verificationToken VerificationToken) SendEmail(tx *db.Tx) error {
	user, err := FetchUserById(tx, verificationToken.UserId)
	if err != nil {
		return err
	}

	err = email.SendEmail(user.Email, "Verify Studier Account", "Hello "+user.Username+",\n\n"+
		"Use this code to verify your Studier account:\n"+verificationToken.Token+"\n\n"+
		"If you did not sign up for a Studier account, you can ignore this email.\n\n"+
		"Thank you.")

	return err
}

func (verificationToken VerificationToken) Delete(tx *db.Tx) error {
	statement := `DELETE FROM verification_tokens WHERE token = $1`
	_, err := tx.Tx.Exec(statement, verificationToken.Token)
	return err
}
