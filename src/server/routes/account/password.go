package account

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type changePasswordRequest struct {
	Ticket      string `json:"ticket"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func ChangePassword(ctx *gin.Context) {
	request := changePasswordRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		panic(err)
	}

	ticketLen := len(request.Ticket)
	oldPasswordLen := len(request.OldPassword)
	if (ticketLen == 0 && oldPasswordLen == 0) ||
		(ticketLen > 0 && oldPasswordLen > 0) ||
		(ticketLen > 0 && ticketLen != 64) ||
		(oldPasswordLen > 0 && (oldPasswordLen < 8 || oldPasswordLen > 2048)) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	if len(request.NewPassword) < 8 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "New password must be at least 8 characters long.",
		})
		return
	} else if len(request.NewPassword) > 2048 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "New password must be less than 2048 characters long.",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	var user *models.User
	if ticketLen > 0 {
		forgotPasswordTicket, err := models.FetchForgotPasswordTicket(tx, request.Ticket)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			panic(err)
		} else if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ticket."})
			return
		}

		err = forgotPasswordTicket.Delete(tx)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		user, err = models.FetchUserById(tx, forgotPasswordTicket.UserId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	} else {
		maybeUser, exists := ctx.Get("user")
		if !exists {
			tx.Rollback()
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You are not authenticated."})
			return
		}

		user = maybeUser.(*models.User)

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.OldPassword))
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			tx.Rollback()
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect old password."})
			return
		} else if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = models.DeleteAllAuthTokensByUserId(tx, user.UserId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = user.SetPassword(tx, string(hashedNewPassword))
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
