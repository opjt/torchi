package handler

import (
	"context"
	"net/http"
	"ohp/internal/api/wrapper"
	"ohp/internal/domain/endpoint"
	"ohp/internal/pkg/log"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type EndpointHandler struct {
	log     *log.Logger
	service *endpoint.EndpointService
}

func NewEndpointHandler(log *log.Logger, service *endpoint.EndpointService) *EndpointHandler {
	return &EndpointHandler{
		log:     log,
		service: service,
	}
}
func (h *EndpointHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", wrapper.WrapJson(h.Add, h.log.Error, wrapper.RespondJSON))
	r.Get("/", h.GetList)
	return r
}

type reqAddEndpoint struct {
	ServiceName string `json:"serviceName"`
}

func (h *EndpointHandler) Add(ctx context.Context, req reqAddEndpoint) (interface{}, error) {

	h.log.Info("req", "sub", req)
	if err := h.service.Add(ctx, req.ServiceName); err != nil {
		return nil, err
	}

	return "success", nil
}

type resListEndpoint struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Endpoint string    `json:"endpoint"`
	Active   bool      `json:"active"`
}

func (h *EndpointHandler) GetList(w http.ResponseWriter, r *http.Request) {

	endpoints, err := h.service.List(r.Context())
	if err != nil {
		wrapper.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	var result []resListEndpoint
	for _, endpoint := range endpoints {
		result = append(result, resListEndpoint{
			ID:       endpoint.ID,
			Name:     endpoint.Name,
			Endpoint: endpoint.Endpoint,
			Active:   true,
		})
	}

	wrapper.RespondJSON(w, http.StatusOK, result)
}
