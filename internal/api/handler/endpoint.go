package handler

import (
	"context"
	"net/http"
	"torchi/internal/api/wrapper"
	"torchi/internal/domain/endpoint"
	"torchi/internal/pkg/log"

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
	r.Delete("/{token}", wrapper.WrapJson(h.Delete, h.log.Error, wrapper.RespondJSON))
	r.Post("/{token}/mute", wrapper.WrapJson(h.Mute, h.log.Error, wrapper.RespondJSON))
	r.Delete("/{token}/mute", wrapper.WrapJson(h.Unmute, h.log.Error, wrapper.RespondJSON))
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

	return nil, nil
}

type resListEndpoint struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Token  string    `json:"token"`
	Active bool      `json:"active"`
}

func (h *EndpointHandler) GetList(w http.ResponseWriter, r *http.Request) {

	endpoints, err := h.service.List(r.Context())
	if err != nil {
		wrapper.RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	result := make([]resListEndpoint, 0, len(endpoints))
	for _, endpoint := range endpoints {
		result = append(result, resListEndpoint{
			ID:     endpoint.ID,
			Name:   endpoint.Name,
			Token:  endpoint.Token,
			Active: endpoint.NotificationEnable,
		})
	}

	wrapper.RespondJSON(w, http.StatusOK, result)
}

func (h *EndpointHandler) Delete(ctx context.Context, _ interface{}) (interface{}, error) {
	token := chi.URLParamFromCtx(ctx, "token")
	if err := h.service.Remove(ctx, token); err != nil {
		return nil, err
	}

	return nil, nil

}

func (h *EndpointHandler) Mute(ctx context.Context, _ interface{}) (interface{}, error) {
	token := chi.URLParamFromCtx(ctx, "token")

	if err := h.service.UpdateMute(ctx, token, false); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *EndpointHandler) Unmute(ctx context.Context, _ interface{}) (interface{}, error) {
	token := chi.URLParamFromCtx(ctx, "token")

	if err := h.service.UpdateMute(ctx, token, true); err != nil {
		return nil, err
	}

	return nil, nil
}
