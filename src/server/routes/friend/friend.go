package friend

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

func DeleteFriend(ctx *gin.Context) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	friend, err := models.FetchFriendByFriendId(tx, ctx.Param("friendId"))
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Friend does not exist."})
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	user := ctx.MustGet("user").(*models.User)

	if friend.ToId != user.UserId && friend.FromId != user.UserId {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own friends."})
		return
	}

	err = friend.Delete(tx)
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

func AcceptFriend(ctx *gin.Context) {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	friend, err := models.FetchFriendByFriendId(tx, ctx.Param("friendId"))
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Friend does not exist."})
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	user := ctx.MustGet("user").(*models.User)

	if friend.ToId != user.UserId && friend.FromId != user.UserId {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own friends."})
		return
	} else if friend.Accepted {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Friend is already accepted."})
		return
	}

	err = friend.SetAccepted(tx, true)
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
