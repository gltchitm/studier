package models

import (
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/id"
)

type Flashcard struct {
	FlashcardId string
	DeckId      string
	Index       int
	Term        string
	Definition  string
}

func CreateFlashcard(tx *db.Tx, deckId string, index int, term, definition string) (*Flashcard, error) {
	flashcardId, err := id.GenerateId(id.IdTypeRegular)
	if err != nil {
		return nil, err
	}

	flashcard := Flashcard{FlashcardId: flashcardId}

	statement := `INSERT INTO flashcards (flashcard_id, deck_id, index, term, definition) VALUES
		($1, $2, $3, $4, $5)`
	_, err = tx.Tx.Exec(statement, flashcardId, deckId, index, term, definition)
	if err != nil {
		return nil, err
	}

	return &flashcard, err
}

func FetchFlashcard(tx *db.Tx, flashcardId string) (*Flashcard, error) {
	flashcard := Flashcard{FlashcardId: flashcardId}
	statement := `SELECT deck_id, index, term, definition FROM flashcards WHERE flashcard_id = $1`
	err := tx.Tx.QueryRow(statement, flashcardId).Scan(&flashcard.DeckId, &flashcard.Index,
		&flashcard.Term, &flashcard.Definition)
	return &flashcard, err
}

func FetchFlashcardsByDeckId(tx *db.Tx, deckId string) (*[]Flashcard, error) {
	statement := `SELECT flashcard_id, index, term, definition FROM flashcards
		WHERE deck_id = $1 ORDER BY index ASC`
	rows, err := tx.Tx.Query(statement, deckId)
	if err != nil {
		return nil, err
	}

	var flashcards []Flashcard
	for rows.Next() {
		flashcard := Flashcard{DeckId: deckId}
		err = rows.Scan(&flashcard.FlashcardId, &flashcard.Index, &flashcard.Term, &flashcard.Definition)
		if err != nil {
			return nil, err
		}
		flashcards = append(flashcards, flashcard)
	}

	return &flashcards, nil
}

func FetchFlashcardByDeckIdAndIndex(tx *db.Tx, deckId string, index int) (*Flashcard, error) {
	flashcard := Flashcard{DeckId: deckId, Index: index}
	statement := `SELECT flashcard_id, term, definition FROM flashcards WHERE deck_id = $1 AND index = $2`
	err := tx.Tx.QueryRow(statement, deckId, index).Scan(&flashcard.FlashcardId, &flashcard.Term,
		&flashcard.Definition)
	return &flashcard, err
}

func CountFlashcardsByDeckId(tx *db.Tx, deckId string) (*int, error) {
	var count int
	statement := `SELECT COUNT(*) FROM flashcards WHERE deck_id = $1`
	err := tx.Tx.QueryRow(statement, deckId).Scan(&count)
	return &count, err
}

func (flashcard *Flashcard) SetTerm(tx *db.Tx, term string) error {
	statement := `UPDATE flashcards SET term = $1 WHERE flashcard_id = $2`
	_, err := tx.Tx.Exec(statement, term, flashcard.FlashcardId)
	flashcard.Term = term
	return err
}

func (flashcard *Flashcard) SetDefinition(tx *db.Tx, definition string) error {
	statement := `UPDATE flashcards SET definition = $1 WHERE flashcard_id = $2`
	_, err := tx.Tx.Exec(statement, definition, flashcard.FlashcardId)
	flashcard.Definition = definition
	return err
}

func (flashcard *Flashcard) SetIndex(tx *db.Tx, index int) error {
	statement := `UPDATE flashcards SET index = $1 WHERE flashcard_id = $2`
	_, err := tx.Tx.Exec(statement, index, flashcard.FlashcardId)
	flashcard.Index = index
	return err
}

func (flashcard *Flashcard) Delete(tx *db.Tx) error {
	statement := `DELETE FROM flashcards WHERE flashcard_id = $1`
	_, err := tx.Tx.Exec(statement, flashcard.FlashcardId)
	return err
}

func DeleteAllFlashcardsByDeckId(tx *db.Tx, deckId string) error {
	statement := `DELETE FROM flashcards WHERE deck_id = $1`
	_, err := tx.Tx.Exec(statement, deckId)
	return err
}
