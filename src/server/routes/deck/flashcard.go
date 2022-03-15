package deck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type newFlashcardRequest struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

func NewFlashcard(ctx *gin.Context) {
	request := newFlashcardRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	} else if len(request.Term) < 3 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Term must be at least 3 characters long.",
		})
		return
	} else if len(request.Term) > 64 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Term must be less than 64 characters long.",
		})
		return
	} else if len(request.Definition) < 3 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Definition must be at least 3 characters long.",
		})
		return
	} else if len(request.Definition) > 64 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Definition must be less than 64 characters long.",
		})
		return
	}

	deckId := ctx.Param("deckId")

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	deck, err := models.FetchDeck(tx, deckId)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Deck does not exist."})
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	if deck.AuthorId != ctx.MustGet("user").(*models.User).UserId {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to edit this deck."})
		return
	}

	flashcardsCount, err := models.CountFlashcardsByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if *flashcardsCount >= 256 {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Deck cannot have more than 256 flashcards."})
		return
	}

	flashcard, err := models.CreateFlashcard(tx, deck.DeckId, *flashcardsCount, request.Term, request.Definition)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"flashcardId": flashcard.FlashcardId})
}
