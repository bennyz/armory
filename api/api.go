package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/middleware"

	"github.com/PUMATeam/catapult/repositories"
	"github.com/PUMATeam/catapult/services"
	"github.com/go-chi/chi"

	"github.com/PUMATeam/catapult/database"
	uuid "github.com/satori/go.uuid"
)

var port int

func New(hs services.Hosts) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)

	hostsEndpoints(r, hs)

	return r
}

func Bootstrap(p int) http.Handler {
	port = p
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	hr := repositories.NewHostsRepository(db)
	hs := services.NewHostsService(hr)
	return New(hs)
}

// Start start the server and listens on the provided port
func Start(h http.Handler) {
	server := http.Server{
		Handler: h,
		Addr:    ":" + strconv.Itoa(port),
	}

	// TODO: add shutdown handling
	server.ListenAndServe()
}

type IDResponse struct {
	ID uuid.UUID `json:"id"`
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case services.ErrNotFound:
		return http.StatusNotFound
	case services.ErrAlreadyExists:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}