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

type GetVarBody struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

type GetVarInput struct {
	VarName string `path:"var_name"`
}

type GetVarResp struct {
	Body GetVarBody `body:""`
}

type SetVarReqBody struct {
	Key   string `json:"key"   required:"true"`
	Value any    `json:"value" required:"true"`
}
type SetVarReq struct {
	Body SetVarReqBody `body:""`
}

type SetVarRespBody struct {
	Message string `json:"message"`
}

type SetVarResp struct {
	Body SetVarRespBody `body:""`
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
		Body: GetVarBody{
			Key:   in.VarName,
			Value: raw,
		},
	}, nil
}

func (h *VarsHandler) SetVar(ctx context.Context, in *SetVarReq) (*SetVarResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if in == nil || in.Body.Key == "" {
		return nil, huma.Error400BadRequest("key is required")
	}

	raw, err := json.Marshal(in.Body.Value)
	if err != nil {
		return nil, huma.Error400BadRequest("invalid value JSON")
	}

	if err := h.svc.SetValue(ctx, in.Body.Key, raw); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, huma.Error404NotFound("variable not found")
		}
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &SetVarResp{
		Body: SetVarRespBody{
			Message: "var successfully updated",
		},
	}, nil
}
