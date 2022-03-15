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

type forgotUsernameRequest struct {
	Email string `json:"email"`
}

func ForgotUsername(ctx *gin.Context) {
	request := forgotUsernameRequest{}

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
		err = email.SendEmail(user.Email, "Studier Username", "Hello,\n\n"+"Your Studier username is:\n"+
			user.Username+"\n\n"+"Thank you.")
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
