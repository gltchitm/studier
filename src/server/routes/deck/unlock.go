package deck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type unlockDeckRequest struct {
	Password string `json:"password"`
}

func UnlockDeck(ctx *gin.Context) {
	request := unlockDeckRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	deckId := ctx.Param("deckId")

	user := ctx.MustGet("user").(*models.User)

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

	author, err := models.FetchUserById(tx, deck.AuthorId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if author.UserId == user.UserId || deck.Access != models.DeckAccessPassword {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You cannot unlock this deck."})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(deck.Password.String), []byte(request.Password))
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password."})
		return
	}

	deckToken, err := models.CreateDeckToken(tx, user.UserId, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"token": deckToken.Token})
}
