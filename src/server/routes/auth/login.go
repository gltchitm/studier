package auth

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	request := loginRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil || len(request.Password) > 2048 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	user, err := models.FetchUserByUsername(tx, request.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	} else if err == nil {
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) == nil {
			authToken, err := models.CreateAuthToken(tx, user.UserId)
			if err != nil {
				tx.Rollback()
				panic(err)
			}

			err = tx.Commit()
			if err != nil {
				panic(err)
			}

			ctx.SetCookie("studier_token", authToken.Token, 9999999999, "/", "", true, true)
			ctx.SetSameSite(http.SameSiteStrictMode)

			ctx.JSON(http.StatusOK, gin.H{"verified": user.Verified})

			return
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username/password."})
}
