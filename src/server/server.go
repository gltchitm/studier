package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gltchtim/studier/server/middleware"
	"github.com/gltchtim/studier/server/routes/account"
	"github.com/gltchtim/studier/server/routes/account/forgot"
	"github.com/gltchtim/studier/server/routes/account/verify"
	"github.com/gltchtim/studier/server/routes/auth"
	"github.com/gltchtim/studier/server/routes/deck"
	"github.com/gltchtim/studier/server/routes/flashcard"
	"github.com/gltchtim/studier/server/routes/friend"
	"github.com/gltchtim/studier/server/routes/user"
)

func main() {
	router := gin.Default()

	accountGroup := router.Group("account")
	accountGroup.GET("", middleware.RequiresAuthMiddleware, account.AccountHandler)
	accountGroup.DELETE("", middleware.RequiresAuthVerifiedMiddleware, account.DeleteAccount)
	accountGroup.POST("signup", middleware.RequiresNoAuthMiddleware, account.Signup)
	accountGroup.POST("email", middleware.RequiresAuthVerifiedMiddleware, account.ChangeEmail)
	accountGroup.POST("password", middleware.MaybeRequiresAuthVerifiedMiddleware, account.ChangePassword)

	accountVerifyGroup := accountGroup.Group("verify")
	accountVerifyGroup.POST("", middleware.RequiresAuthNotVerifiedMiddleware, verify.VerifyAccount)
	accountVerifyGroup.POST("resend", middleware.RequiresAuthNotVerifiedMiddleware,
		verify.ResendVerificationToken)

	accountForgotGroup := accountGroup.Group("forgot")
	accountForgotGroup.POST("username", middleware.RequiresNoAuthMiddleware, forgot.ForgotUsername)

	accountForgotPasswordGroup := accountForgotGroup.Group("password")
	accountForgotPasswordGroup.POST("", middleware.RequiresNoAuthMiddleware, forgot.ForgotPassword)
	accountForgotPasswordGroup.POST("redeem", middleware.RequiresNoAuthMiddleware,
		forgot.RedeemForgotPasswordToken)

	authGroup := router.Group("auth")
	authGroup.POST("login", middleware.RequiresNoAuthMiddleware, auth.Login)
	authGroup.DELETE("logout", middleware.RequiresAuthMiddleware, auth.Logout)

	deckGroup := router.Group("deck")
	deckGroup.GET(":deckId", middleware.RequiresAuthVerifiedMiddleware, deck.Deck)
	deckGroup.DELETE(":deckId", middleware.RequiresAuthVerifiedMiddleware, deck.DeleteDeck)
	deckGroup.POST(":deckId/flashcard", middleware.RequiresAuthMiddleware, deck.NewFlashcard)
	deckGroup.GET("list/all", middleware.RequiresAuthVerifiedMiddleware, deck.ListAllDecks)
	deckGroup.GET("list/pinned", middleware.RequiresAuthVerifiedMiddleware, deck.ListPinnedDecks)
	deckGroup.POST("new", middleware.RequiresAuthVerifiedMiddleware, deck.NewDeck)
	deckGroup.POST(":deckId/unlock", middleware.RequiresAuthVerifiedMiddleware, deck.UnlockDeck)
	deckGroup.POST(":deckId/pin", middleware.RequiresAuthVerifiedMiddleware, deck.PinDeck)
	deckGroup.DELETE(":deckId/pin", middleware.RequiresAuthVerifiedMiddleware, deck.UnpinDeck)
	deckGroup.PUT(":deckId/name", middleware.RequiresAuthVerifiedMiddleware, deck.ChangeDeckName)
	deckGroup.PUT(":deckId/description", middleware.RequiresAuthVerifiedMiddleware, deck.ChangeDeckDescription)
	deckGroup.PUT(":deckId/access", middleware.RequiresAuthVerifiedMiddleware, deck.ChangeDeckAccess)

	flashcardGroup := router.Group("flashcard")
	flashcardGroup.PUT(":flashcardId", middleware.RequiresAuthMiddleware, flashcard.EditFlashcard)
	flashcardGroup.DELETE(":flashcardId", middleware.RequiresAuthMiddleware, flashcard.DeleteFlashcard)
	flashcardGroup.POST(":flashcardId/move", middleware.RequiresAuthMiddleware, flashcard.MoveFlashcard)

	userGroup := router.Group("user")
	userGroup.GET(":userId", middleware.RequiresAuthVerifiedMiddleware, user.User)
	userGroup.POST(":userId/friend", middleware.RequiresAuthVerifiedMiddleware, user.SendFriendRequest)

	friendGroup := router.Group("friend")
	friendGroup.GET("list/incoming", middleware.RequiresAuthVerifiedMiddleware,
		friend.ListIncomingFriendRequests)
	friendGroup.GET("list/outgoing", middleware.RequiresAuthVerifiedMiddleware,
		friend.ListOutgoingFriendRequests)
	friendGroup.GET("list/accepted", middleware.RequiresAuthVerifiedMiddleware,
		friend.ListAcceptedFriends)
	friendGroup.DELETE(":friendId", middleware.RequiresAuthVerifiedMiddleware,
		friend.DeleteFriend)
	friendGroup.POST(":friendId/accept", middleware.RequiresAuthVerifiedMiddleware,
		friend.AcceptFriend)

	router.Run()
}
