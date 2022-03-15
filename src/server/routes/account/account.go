package account

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type deleteAccountRequest struct {
	Password string `json:"password"`
}

func AccountHandler(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	friendRequests, err := models.CountUnacceptedFriendsByToId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{
		"userId":         user.UserId,
		"username":       user.Username,
		"email":          user.Email,
		"verified":       user.Verified,
		"friendRequests": friendRequests,
	})
}

func DeleteAccount(ctx *gin.Context) {
	request := &deleteAccountRequest{}

	err := ctx.ShouldBindJSON(request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
	}

	user := ctx.MustGet("user").(*models.User)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password."})
		return
	} else if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	err = models.DeleteAllAuthTokensByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllDeckTokensByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllPinnedDecksByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	decks, err := models.FetchDecksByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, deck := range *decks {
		err = models.DeleteAllPinnedDecksByDeckId(tx, deck.DeckId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = models.DeleteAllFlashcardsByDeckId(tx, deck.DeckId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = deck.Delete(tx)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = models.DeleteAllFriendsByFromId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllFriendsByToId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = user.Delete(tx)
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
