package user

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type deckInfo struct {
	Id   string `json:"userId"`
	Name string `json:"name"`
}

type userResponse struct {
	Username    string     `json:"username"`
	Decks       []deckInfo `json:"decks"`
	FriendState string     `json:"friendState"`
	FriendId    string     `json:"friendId,omitempty"`
}

const (
	friendStateUnfriended        = "unfriended"
	friendStateFriended          = "friended"
	friendStateRequestedByFriend = "requestedByFriend"
	friendStateRequested         = "requested"
)

func User(ctx *gin.Context) {
	userId := ctx.Param("userId")

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	user, err := models.FetchUserById(tx, userId)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User does not exist."})
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	requestingUser := ctx.MustGet("user").(*models.User)

	friendState := ""

	friend, err := models.FetchFriendByFromIdAndToId(tx, user.UserId, requestingUser.UserId)
	if errors.Is(err, sql.ErrNoRows) {
		friend, err = models.FetchFriendByFromIdAndToId(tx, requestingUser.UserId, user.UserId)
		if errors.Is(err, sql.ErrNoRows) {
			friendState = friendStateUnfriended
		} else if err != nil {
			tx.Rollback()
			panic(err)
		}
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	var friendId string

	if friendState == "" {
		friendId = friend.FriendId

		if friend.Accepted {
			friendState = friendStateFriended
		} else if friend.FromId == user.UserId {
			friendState = friendStateRequestedByFriend
		} else {
			friendState = friendStateRequested
		}
	}

	decks := []deckInfo{}

	rawDecks, err := models.FetchDecksByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, rawDeck := range *rawDecks {
		if requestingUser.UserId == rawDeck.AuthorId || rawDeck.Access == models.DeckAccessEveryone ||
			(rawDeck.Access == models.DeckAccessFriends && friendState == friendStateFriended) {
			decks = append(decks, deckInfo{Id: rawDeck.DeckId, Name: rawDeck.Name})
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, userResponse{
		Username:    user.Username,
		Decks:       decks,
		FriendState: friendState,
		FriendId:    friendId,
	})
}
