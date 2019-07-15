package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/PUMATeam/catapult/services"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func hostsEndpoints(r *chi.Mux, hs services.Hosts) {
	addHostHandler := httptransport.NewServer(
		addHostEndpoint(hs),
		decodeAddHostReq,
		encodeResponse,
	)
	r.Method(http.MethodPost, "/hosts", addHostHandler)

	listHostsHandler := httptransport.NewServer(
		hostsEndpoint(hs),
		httptransport.NopRequestDecoder,
		encodeResponse,
	)
	r.Method(http.MethodGet, "/hosts", listHostsHandler)
}

func addHostEndpoint(svc services.Hosts) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(services.NewHost)
		id, err := svc.AddHost(ctx, req)
		return IDResponse{ID: id}, err
	}
}

func hostsEndpoint(hs services.Hosts) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		hosts, err := hs.ListHosts(ctx)
		return hosts, err
	}
}

func decodeAddHostReq(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var host services.NewHost
	log.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&host)
	return host, err
}