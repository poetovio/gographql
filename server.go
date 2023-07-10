package main

import (
	"log"
	"net/http"
	"os"

	"go-graphql-mongodb-api/database"
	"go-graphql-mongodb-api/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "8080"
const mongoString = "mongodb+srv://admin:admin@cluster0.nikbntq.mongodb.net/?retryWrites=true&w=majority"

func main() {
	database.Connect(mongoString)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", c.Handler(srv))
	http.Handle("/playground", playground.Handler("GraphQL playground", "/graphql"))

	log.Printf("OK -> server is running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
