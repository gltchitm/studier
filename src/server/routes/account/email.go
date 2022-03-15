package account

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/email"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type changeEmailRequest struct {
	NewEmail string `json:"newEmail"`
	Password string `json:"password"`
}

func ChangeEmail(ctx *gin.Context) {
	request := changeEmailRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	if !email.IsValidEmailAddress(request.NewEmail) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email is malformed."})
		return
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

	_, err = models.FetchUserByEmail(tx, request.NewEmail)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	} else if err == nil {
		tx.Rollback()
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email is taken."})
		return
	}

	err = user.SetEmail(tx, request.NewEmail)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = user.SetVerified(tx, false)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = models.DeleteAllAuthTokensByUserId(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	verificationToken, err := models.CreateVerificationToken(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = verificationToken.SendEmail(tx)
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
