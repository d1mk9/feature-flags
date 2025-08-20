package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"feature-flags/pkg/repository"
	"feature-flags/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

type VarsHandler struct {
	svc service.Vars
}

func NewVarsHandler(svc service.Vars) *VarsHandler {
	return &VarsHandler{svc: svc}
}

type GetVarInput struct {
	VarName string `path:"var_name"`
}

type GetVarResp struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type SetVarReq struct {
	Key   string          `json:"key"   required:"true"`
	Value json.RawMessage `json:"value" required:"true"`
}

type SetVarResp struct {
	Message string `json:"message"`
}

func (h *VarsHandler) GetVar(ctx context.Context, in *GetVarInput) (*GetVarResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	raw, err := h.svc.GetValue(ctx, in.VarName)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, huma.Error404NotFound("variable not found")
		}
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &GetVarResp{
		Key:   in.VarName,
		Value: raw,
	}, nil
}

func (h *VarsHandler) SetVar(ctx context.Context, in *SetVarReq) (*SetVarResp, error) {
	if in == nil || in.Key == "" {
		return nil, huma.Error400BadRequest("key is required")
	}

	if err := h.svc.SetValue(ctx, in.Key, in.Value); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, huma.Error404NotFound("variable not found")
		}
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &SetVarResp{Message: "var successfully updated"}, nil
}
