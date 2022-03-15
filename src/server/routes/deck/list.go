package deck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type authorInfo struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type deckInfo struct {
	DeckId      string     `json:"deckId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Flashcards  int        `json:"flashcards"`
	Pinned      bool       `json:"pinned,omitempty"`
	Author      authorInfo `json:"author,omitempty"`
}

func ListAllDecks(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	rawDecks, err := models.FetchDecksByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	decks := []deckInfo{}

	for _, deck := range *rawDecks {
		pinned := false
		_, err = models.FetchPinnedDeckByDeckIdAndUserId(tx, deck.DeckId, user.UserId)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			panic(err)
		} else if err == nil {
			pinned = true
		}

		count, err := models.CountFlashcardsByDeckId(tx, deck.DeckId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		decks = append(decks, deckInfo{
			DeckId:      deck.DeckId,
			Name:        deck.Name,
			Description: deck.Description,
			Flashcards:  *count,
			Pinned:      pinned,
		})
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"decks": decks})
}

func ListPinnedDecks(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	decks := []deckInfo{}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	pinnedDecks, err := models.FetchPinnedDecksByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, pinnedDeck := range *pinnedDecks {
		rawDeck, err := models.FetchDeck(tx, pinnedDeck.DeckId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		count, err := models.CountFlashcardsByDeckId(tx, rawDeck.DeckId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		deck := deckInfo{
			DeckId:      rawDeck.DeckId,
			Name:        rawDeck.Name,
			Description: rawDeck.Description,
			Flashcards:  *count,
		}

		if rawDeck.AuthorId != user.UserId {
			author, err := models.FetchUserById(tx, rawDeck.AuthorId)
			if err != nil {
				tx.Rollback()
				panic(err)
			}

			deck.Author = authorInfo{
				UserId:   author.UserId,
				Username: author.Username,
			}
		}

		decks = append(decks, deck)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"decks": decks})
}
