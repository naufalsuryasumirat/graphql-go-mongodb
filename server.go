package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/naufalsuryasumirat/graphql-go-mongodb/graph"
	"github.com/naufalsuryasumirat/graphql-go-mongodb/graph/generated"
)

const defaultPort = "8080"

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(
		generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()
	r.GET("/", playgroundHandler())
	r.POST("/query", graphqlHandler())
	r.Run(":" + port)
}
