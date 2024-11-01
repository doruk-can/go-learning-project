package http

import (
	_ "foover/docs"
	"foover/internal/middleware"
	"foover/internal/service"
	"foover/internal/transport/http/handler"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
)

func NewRouter(
	sessionService service.SessionService,
	voteService service.VoteService,
	aggregationService service.AggregationService,
	productService service.ProductService,
	logger *slog.Logger,
) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.NewLoggingMiddleware(logger))

	// Session endpoints
	router.HandleFunc("/sessions", handler.CreateSessionHandler(sessionService, logger)).Methods("POST")

	// Vote endpoints
	router.HandleFunc("/votes", handler.SaveVoteHandler(voteService, productService, sessionService, logger)).Methods("POST")
	router.HandleFunc("/votes/{session_id}", handler.GetVotesHandler(voteService, logger)).Methods("GET")

	// Aggregation endpoints
	router.HandleFunc("/aggregated-scores", handler.GetAggregatedScoresHandler(aggregationService, logger)).Methods("GET")

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	// Serving Swagger UI
	router.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	return router
}
