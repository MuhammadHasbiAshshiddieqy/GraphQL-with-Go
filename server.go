package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/generated"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/model"
)

var db *gorm.DB;

func initDB() {
    var err error
		
    db, err = gorm.Open("mysql", os.Getenv("DATA_SOURCE"))

    if err != nil {
        fmt.Println(err)
        panic("failed to connect database")
    }

    db.LogMode(true)

		// db.Exec("CREATE DATABASE db_name")
    // db.Exec("USE db_name") // Use if db name is not define in the db connection

    // Migration to create tables for Order and Item schema
    db.AutoMigrate(&model.GqlCustomer{}, &model.GqlCustomerAddresse{})
}

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	initDB()
	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
        DB: db,
    }}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
