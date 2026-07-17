package server

import (
	"context"
	"github/anurag/altar-be/modules/auth"
	"github/anurag/altar-be/system/config"
	"github/anurag/altar-be/system/constants"
	"github/anurag/altar-be/system/database"
	"github/anurag/altar-be/system/search"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/redis/go-redis/v9"
)

type Server struct {
	CONFIG      *config.Config
	PG_CLIENT   *pgxpool.Pool
	RDS_CLIENT  *redis.Client
	ES_CLIENT   *elasticsearch.Client
	N4J_CLIENT   neo4j.Driver
	ROUTER      *gin.Engine
	HTTP_SERVER *http.Server
}

func NewServer() (*Server, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.CONTEXT_TIMEOUT)
	defer cancel()

	pgClient, err := database.NewPostgresClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := database.RunMigrations(cfg); err != nil {
		return nil, err
	}

	rdsClient, err := database.NewRedisClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	esClient, err := search.NewElasticSearchClient(cfg)
	if err != nil {
		return nil, err
	}

	n4jClient, err := database.NewNeo4jClient(ctx, cfg)
	if err != nil {
		return nil, err
	}

	router := gin.Default()

	httpServer := &http.Server{
		Addr:         constants.SRV_ADDRESS,
		Handler:      router,
		ReadTimeout:  constants.SRV_READ_TIME_OUT,
		WriteTimeout: constants.SRV_WRITE_TIME_OUT,
		IdleTimeout:  constants.SRV_IDLE_TIME_OUT,
	}

	return &Server{
		CONFIG:      cfg,
		PG_CLIENT:   pgClient,
		RDS_CLIENT:  rdsClient,
		ES_CLIENT:   esClient,
		N4J_CLIENT:  n4jClient,
		ROUTER:      router,
		HTTP_SERVER: httpServer,
	}, nil
}

func (s *Server) RegisterRoutes(){
	api := s.ROUTER.Group("/api/v1")
	auth.RegisterAuthRoutes(api.Group("/auth"),authHandler)
}

func (s *Server) Run() error {
	return s.HTTP_SERVER.ListenAndServe()
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.CONTEXT_TIMEOUT)
	defer cancel()

	s.HTTP_SERVER.Shutdown(ctx)
	s.ES_CLIENT.Close(ctx)
	s.PG_CLIENT.Close()
	s.RDS_CLIENT.Close()
	s.N4J_CLIENT.Close(ctx)

}
