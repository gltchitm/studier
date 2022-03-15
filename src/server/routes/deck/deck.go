package deck

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type flashcardInfo struct {
	FlashcardId string `json:"flashcardId"`
	Term        string `json:"term"`
	Definition  string `json:"definition"`
}

func Deck(ctx *gin.Context) {
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

	if user.UserId != author.UserId && deck.Access != models.DeckAccessEveryone {
		accessDenied := true

		switch deck.Access {
		case models.DeckAccessFriends:
			_, err := models.FetchAcceptedFriendByFromIdAndToId(tx, author.UserId, user.UserId)
			if errors.Is(err, sql.ErrNoRows) {
				_, err = models.FetchAcceptedFriendByFromIdAndToId(tx, user.UserId, author.UserId)
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					tx.Rollback()
					panic(err)
				} else if err == nil {
					accessDenied = false
				}
			} else if err != nil {
				tx.Rollback()
				panic(err)
			} else {
				accessDenied = false
			}
		case models.DeckAccessPassword:
			deckToken, err := models.FetchDeckToken(tx, ctx.Query("token"))
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				tx.Rollback()
				panic(err)
			} else if deckToken != nil && deckToken.UserId == user.UserId &&
				deckToken.DeckId == deck.DeckId {
				accessDenied = false
			}
		}

		if accessDenied {
			tx.Rollback()
			ctx.JSON(http.StatusForbidden, gin.H{
				"error":      "You do not have access to this deck.",
				"name":       deck.Name,
				"unlockable": deck.Access == models.DeckAccessPassword,
			})
			return
		}
	}

	rawFlashcards, err := models.FetchFlashcardsByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	flashcards := []flashcardInfo{}

	for _, flashcard := range *rawFlashcards {
		flashcards = append(flashcards, flashcardInfo{
			FlashcardId: flashcard.FlashcardId,
			Term:        flashcard.Term,
			Definition:  flashcard.Definition,
		})
	}

	pinned := false
	_, err = models.FetchPinnedDeckByDeckIdAndUserId(tx, deck.DeckId, user.UserId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	} else if err == nil {
		pinned = true
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"name":        deck.Name,
		"description": deck.Description,
		"access":      deck.Access,
		"flashcards":  flashcards,
		"editable":    deck.AuthorId == ctx.MustGet("user").(*models.User).UserId,
		"author": gin.H{
			"userId":   author.UserId,
			"username": author.Username,
		},
		"pinned": pinned,
	})
}

func DeleteDeck(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

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

	if user.UserId != deck.AuthorId {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to delete this deck."})
		return
	}

	err = models.DeleteAllFlashcardsByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllDeckTokensByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllPinnedDecksByDeckId(tx, deck.DeckId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = deck.Delete(tx)
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
