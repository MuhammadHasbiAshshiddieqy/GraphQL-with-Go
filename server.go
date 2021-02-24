package main

import (
	"fmt"
	"log"
	"net/http"
	"context"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph"
	"github.com/MuhammadHasbiAshshiddieqy/GraphQL-with-Go/graph/generated"
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
}

// HeaderMiddleware - add header as context
func HeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := context.Background()
		profileID := r.Header.Get("Context")
		ctx := context.WithValue(c, graph.Key{}, profileID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}


const defaultPort = "3000"

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
	http.Handle("/query", HeaderMiddleware(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
