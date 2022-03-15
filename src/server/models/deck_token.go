package models

import (
	"database/sql"
	"time"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type DeckToken struct {
	Token   string
	UserId  string
	DeckId  string
	Expires time.Time
}

func CreateDeckToken(tx *db.Tx, userId, deckId string) (*DeckToken, error) {
	token, err := id.GenerateId(id.IdTypeToken)
	if err != nil {
		return nil, err
	}

	deckToken := DeckToken{Token: token, UserId: userId, DeckId: deckId}

	statement := `INSERT INTO deck_tokens (token, user_id, deck_id) VALUES ($1, $2, $3) RETURNING expires`
	err = tx.Tx.QueryRow(statement, token, userId, deckId).Scan(&deckToken.Expires)
	if err != nil {
		return nil, err
	}

	return &deckToken, err
}

func FetchDeckToken(tx *db.Tx, token string) (*DeckToken, error) {
	deckToken := DeckToken{Token: token}
	statement := `SELECT user_id, deck_id, expires FROM deck_tokens WHERE token = $1`
	err := tx.Tx.QueryRow(statement, token).Scan(&deckToken.UserId, &deckToken.DeckId, &deckToken.Expires)
	if err != nil {
		return nil, err
	}
	if deckToken.Expires.Before(time.Now()) {
		statement = `DELETE FROM deck_tokens WHERE token = $1`
		_, err = tx.Tx.Exec(statement, token)
		if err != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	}
	return &deckToken, nil
}

func (deckToken *DeckToken) Delete(tx *db.Tx) error {
	statement := `DELETE FROM deck_tokens WHERE token = $1`
	_, err := tx.Tx.Exec(statement, deckToken.Token)
	return err
}

func DeleteAllDeckTokensByDeckId(tx *db.Tx, deckId string) error {
	statement := `DELETE FROM deck_tokens WHERE deck_id = $1`
	_, err := tx.Tx.Exec(statement, deckId)
	return err
}

func DeleteAllDeckTokensByUserId(tx *db.Tx, userId string) error {
	statement := `DELETE FROM deck_tokens WHERE user_id = $1`
	_, err := tx.Tx.Exec(statement, userId)
	return err
}
