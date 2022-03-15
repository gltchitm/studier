package account

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/email"
	"github.com/gltchtim/studier/server/models"
	"golang.org/x/crypto/bcrypt"
)

type signupRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Signup(ctx *gin.Context) {
	usernameRegexp := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9]+$`)
	request := signupRequest{}

	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Unprocessible entity."})
		return
	} else if !email.IsValidEmailAddress(request.Email) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Email is malformed."})
		return
	} else if !usernameRegexp.MatchString(request.Username) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Username must start with an alphabetic character an be alphanumeric.",
		})
		return
	} else if len(request.Username) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Username must be at least 6 characters long.",
		})
		return
	} else if len(request.Username) > 32 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Username must be less than 32 characters long.",
		})
		return
	} else if len(request.Password) < 8 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Password must be at least 8 characters long.",
		})
		return
	} else if len(request.Password) > 2048 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Password must be less than 2048 characters long.",
		})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = models.FetchUserByEmail(tx, request.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	} else if err == nil {
		tx.Rollback()
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email is taken."})
		return
	}

	_, err = models.FetchUserByUsername(tx, request.Username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		panic(err)
	} else if err == nil {
		tx.Rollback()
		ctx.JSON(http.StatusConflict, gin.H{"error": "Username is taken."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		panic(err)
	}

	user, err := models.CreateUser(tx, request.Email, request.Username, string(hashedPassword))
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
