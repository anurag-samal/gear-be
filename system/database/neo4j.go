package database

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github/anurag/altar-be/system/config"
)

func NewNeo4jClient(ctx context.Context, cfg *config.Config) (neo4j.Driver, error) {

	driver, err := neo4j.NewDriver(cfg.N4J_URL, neo4j.BasicAuth(cfg.N4J_USER, cfg.N4J_PASSWORD, ""))
	if err != nil {
		return nil, err
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		return nil, err
	}
	return driver, nil
}