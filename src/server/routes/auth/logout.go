package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

func Logout(ctx *gin.Context) {
	authToken := ctx.MustGet("auth_token").(*models.AuthToken)

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	err = authToken.Delete(tx)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.SetCookie("studier_token", "", -1, "/", "", true, true)
	ctx.SetSameSite(http.SameSiteStrictMode)

	ctx.JSON(http.StatusOK, gin.H{})
}
