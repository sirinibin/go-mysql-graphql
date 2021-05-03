package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gitlab.com/sirinibin/go-mysql-graphql/config"
	"gitlab.com/sirinibin/go-mysql-graphql/graph"
	"gitlab.com/sirinibin/go-mysql-graphql/graph/generated"
	"gitlab.com/sirinibin/go-mysql-graphql/graph/model"
)

const defaultPort = "8084"

func main() {

	config.InitMysql()
	config.InitRedis()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	c := generated.Config{Resolvers: &graph.Resolver{}}

	c.Directives.IsAuthenticated = func(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {

		UserID := ctx.Value("UserID")
		if UserID != nil {
			return next(ctx)
		} else {
			return nil, errors.New("Unauthorised")
		}
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", AuthHandler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenClaims, err := model.AuthenticateByAccessToken(r)
		if err != nil {
			next.ServeHTTP(w, r)
		} else {
			ctx := context.WithValue(r.Context(), "UserID", tokenClaims.UserID)
			ctx = context.WithValue(ctx, "AccessUUID", tokenClaims.AccessUUID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
