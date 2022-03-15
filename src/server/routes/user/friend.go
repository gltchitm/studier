package user

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

func SendFriendRequest(ctx *gin.Context) {
	from := ctx.MustGet("user").(*models.User)

	toId := ctx.Param("userId")

	if from.UserId == toId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You cannot send a friend request to yourself."})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	to, err := models.FetchUserById(tx, toId)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User does nost exist."})
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	friend, err := models.FetchFriendByFromIdAndToId(tx, from.UserId, to.UserId)
	if !errors.Is(err, sql.ErrNoRows) && err != nil {
		tx.Rollback()
		panic(err)
	} else if err != nil {
		friend, err = models.FetchFriendByFromIdAndToId(tx, to.UserId, from.UserId)

		if !errors.Is(err, sql.ErrNoRows) && err != nil {
			tx.Rollback()
			panic(err)
		} else if err != nil {
			friend, err = models.NewFriend(tx, from.UserId, to.UserId)
			if err != nil {
				tx.Rollback()
				panic(err)
			}

			err = tx.Commit()
			if err != nil {
				panic(err)
			}

			ctx.JSON(http.StatusOK, gin.H{"friendId": friend.FriendId})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	if friend.Accepted {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You are already friends with this user."})
	} else if friend.FromId == from.UserId {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You already sent a friend request to this user."})
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "This user has already sent you a friend request."})
	}
}
