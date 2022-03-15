package verify

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

type verifyAccountRequest struct {
	Token string `json:"verificationCode"`
}

func VerifyAccount(ctx *gin.Context) {
	request := verifyAccountRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	user := ctx.MustGet("user").(*models.User)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	verificationToken, err := models.FetchVerificationToken(tx, user.UserId)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	if request.Token != verificationToken.Token {
		tx.Rollback()
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Incorrect verification token."})
		return
	}

	err = verificationToken.Delete(tx)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = user.SetVerified(tx, true)
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
