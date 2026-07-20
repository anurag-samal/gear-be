package auth

import (
	"github.com/gin-gonic/gin"

	"github/anurag/altar-be/system/middlewares"
	"github/anurag/altar-be/system/packages"
)

func RegisterRoutes(handler *AuthHandler, router *gin.RouterGroup, jwt *pkg.JWTManager) {

	public := router.Group("")

	private := router.Group("")
	private.Use(middlewares.AuthMiddleware(jwt))

	// Email Authentication
	public.POST("/register", handler.Register)
	public.POST("/login", handler.Login)

	// Google OAuth
	public.GET("/google/login", handler.GoogleLogin)
	public.GET("/google/callback", handler.GoogleCallback)

	// Session
	public.POST("/refresh", handler.Refresh)
	public.POST("/logout", handler.Logout)

	// Current User
	private.GET("/profile", handler.Profile)
}