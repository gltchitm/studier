package flashcard

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type moveFlashcardRequest struct {
	Direction string `json:"direction"`
}

const (
	moveFlashcardDirectionUp   = "up"
	moveFlashcardDirectionDown = "down"
)

func MoveFlashcard(ctx *gin.Context) {
	request := moveFlashcardRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	} else if request.Direction != moveFlashcardDirectionUp &&
		request.Direction != moveFlashcardDirectionDown {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid direction."})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	flashcard, err := models.FetchFlashcard(tx, ctx.Param("flashcardId"))
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Flashcard does not exist."})
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	deck, err := models.FetchDeckByFlashcardId(tx, flashcard.FlashcardId)
	if err != nil {
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

	if request.Direction == moveFlashcardDirectionUp {
		if flashcard.Index == 0 {
			tx.Rollback()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Flashcard is already first."})
			return
		}

		preceedingFlashcard, err := models.FetchFlashcardByDeckIdAndIndex(tx, deck.DeckId, flashcard.Index-1)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = flashcard.SetIndex(tx, flashcard.Index-1)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = preceedingFlashcard.SetIndex(tx, preceedingFlashcard.Index+1)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	} else if request.Direction == moveFlashcardDirectionDown {
		if flashcard.Index+1 >= *flashcardsCount {
			tx.Rollback()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Flashcard is already last."})
			return
		}

		nextFlashcard, err := models.FetchFlashcardByDeckIdAndIndex(tx, deck.DeckId, flashcard.Index+1)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = flashcard.SetIndex(tx, flashcard.Index+1)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = nextFlashcard.SetIndex(tx, nextFlashcard.Index-1)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
