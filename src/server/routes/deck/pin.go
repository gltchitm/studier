package deck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

func PinDeck(ctx *gin.Context) {
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

	user := ctx.MustGet("user").(*models.User)

	if deck.AuthorId != user.UserId {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only pin your own decks."})
		return
	}

	_, err = models.FetchPinnedDeckByDeckIdAndUserId(tx, deck.DeckId, user.UserId)
	if err == nil {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Deck is already pinned."})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	}

	_, err = models.CreatePinnedDeck(tx, deck.DeckId, user.UserId)
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

func UnpinDeck(ctx *gin.Context) {
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

	user := ctx.MustGet("user").(*models.User)

	if deck.AuthorId != user.UserId {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only pin your own decks."})
		return
	}

	pinnedDeck, err := models.FetchPinnedDeckByDeckIdAndUserId(tx, deck.DeckId, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = pinnedDeck.Delete(tx)
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
