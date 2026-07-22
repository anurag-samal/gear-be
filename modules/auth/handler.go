package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github/anurag/altar-be/system/config"
	"github/anurag/altar-be/modules/auth/utils"
	"github/anurag/altar-be/system/middlewares"
	"net/http"
)

type AuthHandler struct {
	service *AuthService
	config  *config.Config
}

func NewHandler( service *AuthService, cfg *config.Config) *AuthHandler {

	return &AuthHandler{
		service: service,
		config:  cfg,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req auth.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, refreshToken, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setRefreshCookie(c, refreshToken)

	c.JSON(http.StatusCreated, res)
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req auth.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, refreshToken, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setRefreshCookie(c, refreshToken)

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) GoogleLogin(c *gin.Context) {

    url, err := h.service.GoogleLogin(c.Request.Context())
    if err != nil {
        h.handleError(c, err)
        return
    }

    c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {

	code := c.Query("code")
	state := c.Query("state")

	_, refreshToken, err := h.service.GoogleCallback(
		c.Request.Context(),
		code,
		state,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setRefreshCookie(c, refreshToken)

	c.Redirect(
		http.StatusTemporaryRedirect,
		h.config.FrontendURL+"/dashboard",
	)
}

func (h *AuthHandler) Refresh(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "missing refresh token",
		})
		return
	}

	res, newRefreshToken, err := h.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.setRefreshCookie(c, newRefreshToken)

	c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		_ = h.service.Logout(c.Request.Context(), refreshToken)
	}

	h.clearRefreshCookie(c)

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) Profile(c *gin.Context) {

	claims := middlewares.GetClaims(c)
	if claims == nil {
		h.handleError(c, auth.ErrUnauthorized)
		return
	}

	res, err := h.service.Profile(
		c.Request.Context(),
		claims.UserID,
	)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

// ----------------------------------------------------------------------
// Cookie Helpers
// ----------------------------------------------------------------------

func (h *AuthHandler) setRefreshCookie(c *gin.Context, token string) {

	c.SetCookie(
		"refresh_token",
		token,
		int(auth.JWT_REFRESH_TOKEN_TTL.Seconds()),
		"/auth",
		"",
		h.config.CookieSecure,
		true,
	)
}

func (h *AuthHandler) clearRefreshCookie(c *gin.Context) {

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/auth",
		"",
		h.config.CookieSecure,
		true,
	)
}

// ----------------------------------------------------------------------
// Error Handling
// ----------------------------------------------------------------------

func (h *AuthHandler) handleError(c *gin.Context, err error) {

	switch {

	case errors.Is(err, auth.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, auth.ErrUserAlreadyExists):
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, auth.ErrRefreshTokenInvalid):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, auth.ErrRefreshTokenExpired):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

	case errors.Is(err, auth.ErrRefreshTokenRevoked):
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}
}
