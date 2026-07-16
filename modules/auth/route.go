package auth

import (
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, handler *AuthHandler) {
	public := router
	protected := router.Group("")
	protected.Use(AuthMiddleware())
	
	// Email Authentication
	public.POST("/register", handler.Register)
	public.POST("/login", handler.Login)

	// Google OIDC + PKCE
	public.GET("/google/login", handler.GoogleLogin)
	public.GET("/google/callback", handler.GoogleCallback)

	// Session Management
	public.POST("/refresh", handler.Refresh)
	protected.POST("/logout", handler.Logout)

	// Current User
	protected.GET("/profile", handler.Profile)
}
