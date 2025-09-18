package handler

import (
	"github.com/toastsandwich/kvstore/pkg/kvstore"
)

type Body struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Handler struct {
	grpcEndpoint string
	s            *kvstore.Store
}

func New(store *kvstore.Store) *Handler {
	return &Handler{
		s: store,
	}
}

func (h *Handler) SetGrpcEndpoint(s string) *Handler {
	h.grpcEndpoint = s
	return h
}
