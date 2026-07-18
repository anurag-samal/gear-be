package auth

import (
	"github/anurag/altar-be/system/config"
	pkg "github/anurag/altar-be/system/packages"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func BuildAuthModule(router *gin.RouterGroup, db *pgxpool.Pool, redis *redis.Client, cfg *config.Config, jwt *pkg.JWTManager, google *pkg.GoogleProvider) {
	repo := NewRepository()
	service := NewService(repo, db, redis, cfg, jwt, google)
	handler := NewHandler(service, cfg)
	RegisterRoutes(handler, router, jwt)
}
