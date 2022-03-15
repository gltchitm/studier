package flashcard

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type editFlashcardRequest struct {
	Term       string `json:"term"`
	Definition string `json:"definition"`
}

func EditFlashcard(ctx *gin.Context) {
	request := editFlashcardRequest{}

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

	err = flashcard.SetTerm(tx, request.Term)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = flashcard.SetDefinition(tx, request.Definition)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func DeleteFlashcard(ctx *gin.Context) {
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
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to delete this deck."})
		return
	}

	count, err := models.CountFlashcardsByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if *count <= 3 {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Decks must have at least 3 flashcards."})
		return
	}

	err = flashcard.Delete(tx)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	flashcards, err := models.FetchFlashcardsByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for i, otherFlashcard := range (*flashcards)[flashcard.Index:] {
		err = otherFlashcard.SetIndex(tx, i+flashcard.Index)
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
