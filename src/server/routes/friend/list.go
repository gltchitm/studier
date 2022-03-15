package friend

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type userInfo struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type friendInfo struct {
	FriendId  string    `json:"friendId"`
	Timestamp int64     `json:"timestamp,omitempty"`
	From      *userInfo `json:"from,omitempty"`
	To        *userInfo `json:"to,omitempty"`
	Friend    *userInfo `json:"friend,omitempty"`
}

func ListIncomingFriendRequests(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	rawFriends, err := models.FetchFriendsByToId(tx, false, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	friends := []friendInfo{}

	for _, rawFriend := range *rawFriends {
		from, err := models.FetchUserById(tx, rawFriend.FromId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		friends = append(friends, friendInfo{
			FriendId:  rawFriend.FriendId,
			Timestamp: rawFriend.Timestamp.UnixMilli(),
			From: &userInfo{
				UserId:   from.UserId,
				Username: from.Username,
			},
		})
	}

	tx.Commit()

	ctx.JSON(http.StatusOK, gin.H{"friendRequests": friends})
}

func ListOutgoingFriendRequests(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	rawFriends, err := models.FetchFriendsByFromId(tx, false, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	friends := []friendInfo{}

	for _, rawFriend := range *rawFriends {
		to, err := models.FetchUserById(tx, rawFriend.ToId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		friends = append(friends, friendInfo{
			FriendId: rawFriend.FriendId,
			To: &userInfo{
				UserId:   to.UserId,
				Username: to.Username,
			},
		})
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"friendRequests": friends})
}

func ListAcceptedFriends(ctx *gin.Context) {
	user := ctx.MustGet("user").(*models.User)

	friends := []friendInfo{}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	rawFriends, err := models.FetchFriendsByFromId(tx, true, user.UserId)
	if err != nil {
		panic(err)
	}

	for _, rawFriend := range *rawFriends {
		to, err := models.FetchUserById(tx, rawFriend.ToId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		friends = append(friends, friendInfo{
			FriendId: rawFriend.FriendId,
			Friend: &userInfo{
				UserId:   to.UserId,
				Username: to.Username,
			},
		})
	}

	rawFriends, err = models.FetchFriendsByToId(tx, true, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	for _, rawFriend := range *rawFriends {
		from, err := models.FetchUserById(tx, rawFriend.FromId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		friends = append(friends, friendInfo{
			FriendId:  rawFriend.FriendId,
			Timestamp: rawFriend.Timestamp.UnixMilli(),
			Friend: &userInfo{
				UserId:   from.UserId,
				Username: from.Username,
			},
		})
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"friends": friends})
}
