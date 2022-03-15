package forgot

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/email"
	"github.com/gltchtim/studier/server/models"
)

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func ForgotPassword(ctx *gin.Context) {
	request := forgotPasswordRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	if !email.IsValidEmailAddress(request.Email) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email is malformed."})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	user, err := models.FetchUserByEmail(tx, request.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	}

	if user != nil {
		forgotPasswordToken, err := models.CreateForgotPasswordToken(tx, user.UserId)
		if err != nil {
			tx.Rollback()
			panic(err)
		}

		err = forgotPasswordToken.SendEmail(tx)
		if err != nil {
			tx.Rollback()
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

type redeemForgotPasswordTokenRequest struct {
	Token string `json:"token"`
}

func RedeemForgotPasswordToken(ctx *gin.Context) {
	request := redeemForgotPasswordTokenRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	forgotPasswordToken, err := models.FetchForgotPasswordToken(tx, request.Token)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	} else if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Invalid token."})
		return
	}

	err = forgotPasswordToken.Delete(tx)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	forgotPasswordTicket, err := models.CreateForgotPasswordTicket(tx, forgotPasswordToken.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusOK, gin.H{"ticket": forgotPasswordTicket.Ticket})
}
