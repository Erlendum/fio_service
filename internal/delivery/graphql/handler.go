package graphql

import (
	graph2 "fio_service/internal/delivery/graphql/graph"
	"fio_service/internal/service"
	"fio_service/pkg/logger"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"net/http"
)

type Handler struct {
	srv *handler.Server
}

func NewHandler(services *service.Services, logger *logger.Logger) *Handler {
	srv := handler.NewDefaultServer(graph2.NewExecutableSchema(graph2.Config{Resolvers: &graph2.Resolver{
		Services: *services,
		Logger:   *logger}}))

	return &Handler{srv: srv}

}

func (h *Handler) Init() *handler.Server {
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", h.srv)
	return h.srv
}
