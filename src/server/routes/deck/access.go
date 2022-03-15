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

type changeAccessRequest struct {
	Access   string `json:"access"`
	Password string `json:"password"`
}

func ChangeDeckAccess(ctx *gin.Context) {
	request := &changeAccessRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil || (request.Access != models.DeckAccessEveryone &&
		request.Access != models.DeckAccessFriends &&
		request.Access != models.DeckAccessPassword &&
		request.Access != models.DeckAccessMe) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	} else if request.Access == models.DeckAccessPassword {
		if request.Access == models.DeckAccessPassword && len(request.Password) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Password must be more than 6 characters long.",
			})
			return
		} else if request.Access == models.DeckAccessPassword && len(request.Password) > 64 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Password must be less than 64 characters long.",
			})
			return
		}
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

	err = deck.SetAccess(tx, request.Access)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	password := sql.NullString{}
	if request.Access == models.DeckAccessPassword {
		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			tx.Rollback()
			panic(err)
		}

		password.String = string(hashedPassword)
		password.Valid = true
	}

	err = deck.SetPassword(tx, password)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllDeckTokensByDeckId(tx, deck.DeckId)
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
