package database

import (
    "context"
    "fmt"
    "github.com/neo4j/neo4j-go-driver/v6/neo4j"
)

func main() {
    ctx := context.Background()
    dbUri := "neo4j+s://8c3b2865.databases.neo4j.io"
    dbUser := "8c3b2865"
    dbPassword := "<password>"
    driver, err := neo4j.NewDriver(
        dbUri,
        neo4j.BasicAuth(dbUser, dbPassword, ""))
    if err != nil {
        panic(err)
    }
    defer driver.Close(ctx)

    err = driver.VerifyConnectivity(ctx)
    if err != nil {
        panic(err)
    }
    fmt.Println("Connection established.")
}