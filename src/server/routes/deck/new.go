package deck

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type newDeckRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Access      string `json:"access"`
	Password    string `json:"password"`
	Flashcards  []struct {
		Term       string `json:"term"`
		Definition string `json:"definition"`
	} `json:"flashcards"`
}

func NewDeck(ctx *gin.Context) {
	request := newDeckRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil || (request.Access != models.DeckAccessEveryone &&
		request.Access != models.DeckAccessFriends &&
		request.Access != models.DeckAccessPassword &&
		request.Access != models.DeckAccessMe) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	} else if len(request.Flashcards) < 3 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Deck must have at least 3 flashcards.",
		})
		return
	} else if len(request.Flashcards) > 256 {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Deck cannot have more than 256 flashcards.",
		})
		return
	} else if len(request.Name) < 3 {
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": "Name must be more than 3 characters long.",
		})
		return
	} else if len(request.Name) > 32 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Name must be less than 32 characters long.",
		})
		return
	} else if len(request.Description) > 96 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Description must be less than 96 characters long.",
		})
		return
	} else if request.Access == models.DeckAccessPassword && len(request.Password) < 6 {
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

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	user := ctx.MustGet("user").(*models.User)

	password := sql.NullString{}

	if request.Access == models.DeckAccessPassword {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		password.String = string(hashedPassword)
		password.Valid = true
	}

	deck, err := models.CreateDeck(tx, user.UserId, request.Name, request.Description, request.Access,
		password)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for i, flashcardRequest := range request.Flashcards {
		if len(flashcardRequest.Term) < 3 {
			tx.Rollback()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Flashcard " + fmt.Sprint(i+1) + ": Term must be at least 3 characters long.",
			})
			return
		} else if len(flashcardRequest.Term) > 64 {
			tx.Rollback()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Flashcard " + fmt.Sprint(i+1) + ": Term must be less than 64 characters long.",
			})
			return
		} else if len(flashcardRequest.Definition) < 3 {
			tx.Rollback()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Flashcard " + fmt.Sprint(i+1) + ": Definition must be at least 3 characters long.",
			})
			return
		} else if len(flashcardRequest.Definition) > 64 {
			tx.Rollback()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": "Flashcard " + fmt.Sprint(i+1) + ": Definition must be less than 64 characters long.",
			})
			return
		}

		_, err := models.CreateFlashcard(tx, deck.DeckId, i, flashcardRequest.Term,
			flashcardRequest.Definition)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"deckId": deck.DeckId})
}
