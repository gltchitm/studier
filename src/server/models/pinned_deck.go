package models

import (
	"time"

	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type PinnedDeck struct {
	PinnedDeckId string
	Created      time.Time
	DeckId       string
	UserId       string
}

func CreatePinnedDeck(tx *db.Tx, deckId, userId string) (*PinnedDeck, error) {
	pinnedDeckId, err := id.GenerateId(id.IdTypeRegular)
	if err != nil {
		return nil, err
	}

	pinnedDeck := PinnedDeck{PinnedDeckId: pinnedDeckId, DeckId: deckId, UserId: userId}

	statement := `INSERT INTO pinned_decks (pinned_deck_id, deck_id, user_id) VALUES ($1, $2, $3)`
	_, err = tx.Tx.Exec(statement, pinnedDeckId, deckId, userId)
	if err != nil {
		return nil, err
	}

	return &pinnedDeck, nil
}

func FetchPinnedDeckByDeckIdAndUserId(tx *db.Tx, deckId, userId string) (*PinnedDeck, error) {
	pinnedDeck := PinnedDeck{DeckId: deckId, UserId: userId}
	statement := `SELECT pinned_deck_id, created FROM pinned_decks WHERE deck_id = $1 AND user_id = $2`
	err := tx.Tx.QueryRow(statement, deckId, userId).Scan(&pinnedDeck.PinnedDeckId, &pinnedDeck.Created)
	return &pinnedDeck, err
}

func FetchPinnedDecksByUserId(tx *db.Tx, userId string) (*[]PinnedDeck, error) {
	statement := `SELECT pinned_deck_id, created, deck_id FROM pinned_decks WHERE user_id = $1
		ORDER BY created DESC`
	rows, err := tx.Tx.Query(statement, userId)
	if err != nil {
		return nil, err
	}

	var pinnedDecks []PinnedDeck
	for rows.Next() {
		pinnedDeck := PinnedDeck{UserId: userId}
		err = rows.Scan(&pinnedDeck.PinnedDeckId, &pinnedDeck.Created, &pinnedDeck.DeckId)
		if err != nil {
			return nil, err
		}
		pinnedDecks = append(pinnedDecks, pinnedDeck)
	}

	return &pinnedDecks, nil
}

func DeleteAllPinnedDecksByDeckId(tx *db.Tx, deckId string) error {
	statement := `DELETE FROM pinned_decks WHERE deck_id = $1`
	_, err := tx.Tx.Exec(statement, deckId)
	return err
}

func DeleteAllPinnedDecksByUserId(tx *db.Tx, userId string) error {
	statement := `DELETE FROM pinned_decks WHERE user_id = $1`
	_, err := tx.Tx.Exec(statement, userId)
	return err
}

func CountPinnedDecksByUserId(tx *db.Tx, userId string) (*int, error) {
	var count int
	statement := `SELECT COUNT(*) FROM pinned_decks WHERE user_id = $1`
	err := tx.Tx.QueryRow(statement, userId).Scan(&count)
	return &count, err
}

func (pinnedDeck *PinnedDeck) Delete(tx *db.Tx) error {
	statement := `DELETE FROM pinned_decks WHERE pinned_deck_id = $1`
	_, err := tx.Tx.Exec(statement, pinnedDeck.PinnedDeckId)
	return err
}
