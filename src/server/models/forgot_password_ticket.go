package models

import (
	"database/sql"
	"time"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type forgotPasswordTicket struct {
	Ticket  string
	UserId  string
	Expires time.Time
}

func CreateForgotPasswordTicket(tx *db.Tx, userId string) (*forgotPasswordTicket, error) {
	ticket, err := id.GenerateId(id.IdTypeTicket)
	if err != nil {
		return nil, err
	}

	forgotPasswordTicket := forgotPasswordTicket{Ticket: ticket, UserId: userId}

	statement := `INSERT INTO forgot_password_tickets (ticket, user_id) VALUES ($1, $2) RETURNING expires`
	err = tx.Tx.QueryRow(statement, ticket, userId).Scan(&forgotPasswordTicket.Expires)
	if err != nil {
		return nil, err
	}
	return &forgotPasswordTicket, err
}

func FetchForgotPasswordTicket(tx *db.Tx, ticket string) (*forgotPasswordTicket, error) {
	forgotPasswordTicket := forgotPasswordTicket{Ticket: ticket}
	statement := `SELECT user_id, expires FROM forgot_password_tickets WHERE ticket = $1`
	err := tx.Tx.QueryRow(statement, ticket).Scan(
		&forgotPasswordTicket.UserId,
		&forgotPasswordTicket.Expires,
	)
	if err != nil {
		return nil, err
	}
	if forgotPasswordTicket.Expires.Before(time.Now()) {
		statement = `DELETE FROM forgot_password_tickets WHERE ticket = $1`
		_, err = tx.Tx.Exec(statement, forgotPasswordTicket.Ticket)
		if err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	return &forgotPasswordTicket, nil
}

func (forgotPasswordTicket *forgotPasswordTicket) Delete(tx *db.Tx) error {
	statement := `DELETE FROM forgot_password_tickets WHERE ticket = $1`
	_, err := tx.Tx.Exec(statement, forgotPasswordTicket.Ticket)
	return err
}
