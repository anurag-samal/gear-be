package database

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github/anurag/altar-be/system/config"
)

func NewNeo4jClient(ctx context.Context, cfg *config.Config) (neo4j.Driver, error) {
	driver, err := neo4j.NewDriver(cfg.Neo4j.URL, neo4j.BasicAuth(cfg.Neo4j.User, cfg.Neo4j.Password, ""))
	if err != nil {
		return nil, fmt.Errorf("neo4j driver: %w", err)
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		driver.Close(ctx)
		return nil, fmt.Errorf("neo4j connectivity: %w", err)
	}

	return driver, nil
}
