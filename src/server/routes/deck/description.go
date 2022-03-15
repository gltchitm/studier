package deck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type changeDeckDescriptionRequest struct {
	Description string `json:"description"`
}

func ChangeDeckDescription(ctx *gin.Context) {
	request := changeDeckDescriptionRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	} else if len(request.Description) < 3 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Description must be more than 3 characters long.",
		})
		return
	} else if len(request.Description) > 64 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Description must be less than 64 characters long.",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	deck, err := models.FetchDeck(tx, ctx.Param("deckId"))
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

	err = deck.SetDescription(tx, request.Description)
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
