package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"

	"github/anurag/altar-be/system/config"
	"github/anurag/altar-be/system/constants"
	"github/anurag/altar-be/system/database"
	pkg "github/anurag/altar-be/system/packages"
	"github/anurag/altar-be/system/search"
)

type Server struct {
	CONFIG         *config.Config
	PG_CLIENT      *database.PostgresConfig
	RDS_CLIENT     *database.RedisConfig
	ES_CLIENT      *elasticsearch.Client
	N4J_CLIENT     neo4j.Driver
	ROUTER         *gin.Engine
	HTTP_SERVER    *http.Server
	JWT_MANAGER    *pkg.JWTManager
	GOOGLE_PROVIDER *pkg.GoogleProvider
}

func NewServer() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	cfg.Print()

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}

	slog.Info("connecting to postgres")
	ctx, cancel := context.WithTimeout(context.Background(), constants.CONTEXT_TIMEOUT)
	defer cancel()

	pgClient, err := database.NewPostgresClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	slog.Info("postgres connected")

	if err := database.RunMigrations(cfg); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}
	slog.Info("migrations complete")

	slog.Info("connecting to redis")
	rdsClient, err := database.NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}
	slog.Info("redis connected")

	slog.Info("connecting to elasticsearch")
	esClient, err := search.NewElasticSearchClient(cfg)
	if err != nil {
		slog.Warn("elasticsearch unavailable (skipping)", "error", err)
		esClient = nil
	}

	slog.Info("connecting to neo4j")
	n4jClient, err := database.NewNeo4jClient(ctx, cfg)
	if err != nil {
		slog.Warn("neo4j unavailable (skipping)", "error", err)
		n4jClient = nil
	}

	slog.Info("initializing jwt manager")
	jwtManager, err := pkg.NewJWTManager(cfg.JWT.Secret)
	if err != nil {
		return nil, fmt.Errorf("jwt: %w", err)
	}

	slog.Info("initializing google oauth provider")
	googleProvider, err := pkg.NewGoogleProvider(ctx, cfg.OAuth.GoogleClientID, cfg.OAuth.GoogleClientSecret, cfg.OAuth.GoogleRedirectURL)
	if err != nil {
		return nil, fmt.Errorf("google oauth: %w", err)
	}

	router := gin.Default()

	addr := fmt.Sprintf(":%s", cfg.App.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  constants.SRV_READ_TIME_OUT,
		WriteTimeout: constants.SRV_WRITE_TIME_OUT,
		IdleTimeout:  constants.SRV_IDLE_TIME_OUT,
	}

	return &Server{
		CONFIG:         cfg,
		PG_CLIENT:      pgClient,
		RDS_CLIENT:     rdsClient,
		ES_CLIENT:      esClient,
		N4J_CLIENT:     n4jClient,
		ROUTER:         router,
		HTTP_SERVER:    httpServer,
		JWT_MANAGER:    jwtManager,
		GOOGLE_PROVIDER: googleProvider,
	}, nil
}

func (s *Server) Run() error {
	slog.Info("starting server", "addr", s.HTTP_SERVER.Addr)
	s.BuildModules()
	return s.HTTP_SERVER.ListenAndServe()
}

func (s *Server) Shutdown() {
	slog.Info("shutting down server")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), constants.CONTEXT_TIMEOUT)
	defer shutdownCancel()

	if err := s.HTTP_SERVER.Shutdown(shutdownCtx); err != nil {
		slog.Error("http server shutdown", "error", err)
	} else {
		slog.Info("http server stopped")
	}

	if s.ES_CLIENT != nil {
		s.ES_CLIENT.Close(shutdownCtx)
		slog.Info("elasticsearch client closed")
	}

	s.PG_CLIENT.Disconnect()

	s.RDS_CLIENT.Disconnect()

	if s.N4J_CLIENT != nil {
		s.N4J_CLIENT.Close(shutdownCtx)
		slog.Info("neo4j driver closed")
	}

	slog.Info("shutdown complete")
}

func (s *Server) ShutdownWithTimeout(d time.Duration) {
	done := make(chan struct{})
	go func() {
		s.Shutdown()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(d):
		slog.Error("forced shutdown after timeout", "timeout", d)
	}
}
