package middleware

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/db"
	"github.com/gltchtim/studier/server/models"
)

const (
	requiresAuthVerified = iota
	requiresAuthNotVerified
	requiresAuthVerifiedOrNotVerified
)

func requiresAuthMiddleware(ctx *gin.Context, requirement int, maybe bool) {
	token, err := ctx.Cookie("studier_token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		panic(err)
	} else if err == nil {
		tx, err := db.Begin()
		if err != nil {
			panic(err)
		}

		authToken, err := models.FetchAuthToken(tx, token)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			tx.Rollback()
			panic(err)
		} else if err == nil {
			user, err := models.FetchUserById(tx, authToken.UserId)
			if err != nil {
				tx.Rollback()
				panic(err)
			}

			if !user.Verified && requirement == requiresAuthVerified {
				tx.Rollback()
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Your account is not verified.",
				})
				return
			} else if user.Verified && requirement == requiresAuthNotVerified {
				tx.Rollback()
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error": "Your account is already verified.",
				})
				return
			}

			err = authToken.Extend(tx)
			if err != nil {
				tx.Rollback()
				panic(err)
			}

			ctx.Set("user", user)
			ctx.Set("auth_token", authToken)

			err = tx.Commit()
			if err != nil {
				panic(err)
			}

			return
		}

		tx.Rollback()
	}

	if !maybe {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are not authenticated."})
	}
}

func RequiresAuthMiddleware(ctx *gin.Context) {
	requiresAuthMiddleware(ctx, requiresAuthVerifiedOrNotVerified, false)
}
func RequiresAuthVerifiedMiddleware(ctx *gin.Context) {
	requiresAuthMiddleware(ctx, requiresAuthVerified, false)
}
func RequiresAuthNotVerifiedMiddleware(ctx *gin.Context) {
	requiresAuthMiddleware(ctx, requiresAuthNotVerified, false)
}

func MaybeRequiresAuthVerifiedMiddleware(ctx *gin.Context) {
	requiresAuthMiddleware(ctx, requiresAuthVerified, true)
}

func RequiresNoAuthMiddleware(ctx *gin.Context) {
	token, err := ctx.Cookie("studier_token")
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		panic(err)
	} else if err != nil {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = models.FetchAuthToken(tx, token)
	if errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return
	} else if err != nil {
		tx.Rollback()
		panic(err)
	}

	tx.Rollback()

	ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You are already authenticated."})
}
