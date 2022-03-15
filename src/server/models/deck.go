package models

import (
	"database/sql"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

const (
	DeckAccessEveryone = "everyone"
	DeckAccessFriends  = "friends"
	DeckAccessPassword = "password"
	DeckAccessMe       = "me"
)

type Deck struct {
	DeckId      string
	AuthorId    string
	Name        string
	Description string
	Access      string
	Password    sql.NullString
}

func CreateDeck(
	tx *db.Tx,
	userId,
	name,
	description,
	access string,
	password sql.NullString,
) (*Deck, error) {
	deckId, err := id.GenerateId(id.IdTypeRegular)
	if err != nil {
		return nil, err
	}

	deck := Deck{DeckId: deckId}

	statement := `INSERT INTO decks (deck_id, author_id, name, description, access,
		password) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Tx.Exec(statement, deckId, userId, name, description, access, password)
	if err != nil {
		return nil, err
	}

	return &deck, err
}

func FetchDecksByUserId(tx *db.Tx, authorId string) (*[]Deck, error) {
	statement := `SELECT deck_id, name, description, access, password FROM decks WHERE author_id = $1`
	rows, err := tx.Tx.Query(statement, authorId)
	if err != nil {
		return nil, err
	}
	decks := []Deck{}
	for rows.Next() {
		deck := Deck{AuthorId: authorId}
		err = rows.Scan(&deck.DeckId, &deck.Name, &deck.Description, &deck.Access,
			&deck.Password)
		if err != nil {
			return nil, err
		}
		decks = append(decks, deck)
	}
	return &decks, nil
}

func FetchDeck(tx *db.Tx, deckId string) (*Deck, error) {
	deck := Deck{DeckId: deckId}
	statement := `SELECT author_id, name, description, access, password FROM decks WHERE deck_id = $1`
	err := tx.Tx.QueryRow(statement, deckId).Scan(&deck.AuthorId, &deck.Name, &deck.Description,
		&deck.Access, &deck.Password)
	return &deck, err
}

func FetchDeckByFlashcardId(tx *db.Tx, flashcardId string) (*Deck, error) {
	var deckId string

	statement := `SELECT deck_id FROM flashcards WHERE flashcard_id = $1`
	err := tx.Tx.QueryRow(statement, flashcardId).Scan(&deckId)
	if err != nil {
		return nil, err
	}
	return FetchDeck(tx, deckId)
}

func (deck *Deck) SetName(tx *db.Tx, name string) error {
	statement := `UPDATE decks SET name = $1 WHERE deck_id = $2`
	_, err := tx.Tx.Exec(statement, name, deck.DeckId)
	deck.Name = name
	return err
}

func (deck *Deck) SetAccess(tx *db.Tx, access string) error {
	statement := `UPDATE decks SET access = $1 WHERE deck_id = $2`
	_, err := tx.Tx.Exec(statement, access, deck.DeckId)
	deck.Access = access
	return err
}

func (deck *Deck) SetPassword(tx *db.Tx, password sql.NullString) error {
	statement := `UPDATE decks SET password = $1 WHERE deck_id = $2`
	_, err := tx.Tx.Exec(statement, password, deck.DeckId)
	deck.Password = password
	return err
}

func (deck *Deck) SetDescription(tx *db.Tx, description string) error {
	statement := `UPDATE decks SET description = $1 WHERE deck_id = $2`
	_, err := tx.Tx.Exec(statement, description, deck.DeckId)
	deck.Description = description
	return err
}

func (deck *Deck) Delete(tx *db.Tx) error {
	statement := `DELETE FROM decks WHERE deck_id = $1`
	_, err := tx.Tx.Exec(statement, deck.DeckId)
	return err
}
