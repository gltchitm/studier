package verify

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

func ResendVerificationToken(ctx *gin.Context) {
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
