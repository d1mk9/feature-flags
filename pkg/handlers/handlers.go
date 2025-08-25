package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"feature-flags/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

type FeatureHandler struct {
	svc service.Flags
}

func NewFeatureHandler(svc service.Flags) *FeatureHandler {
	return &FeatureHandler{svc: svc}
}

type GetVarBody struct {
	Key   string               `json:"key"`
	Value service.FeatureValue `json:"value"`
}

type GetVarParams struct {
	VarName string `path:"var_name"`
}

type GetVarResp struct {
	Body GetVarBody `body:""`
}

type SetVarReqBody struct {
	Key   string          `json:"key"   required:"true"`
	Value json.RawMessage `json:"value" required:"true"`
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

func (h *FeatureHandler) GetVar(ctx context.Context, in *GetVarParams) (*GetVarResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	val, err := h.svc.GetValue(ctx, in.VarName)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			return nil, huma.Error404NotFound("variable not found")
		}
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &GetVarResp{
		Body: GetVarBody{
			Key:   in.VarName,
			Value: val,
		},
	}, nil
}

func (h *FeatureHandler) SetVar(ctx context.Context, in *SetVarReq) (*SetVarResp, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if in == nil || in.Body.Key == "" {
		return nil, huma.Error400BadRequest("key is required")
	}

	var fv service.FeatureValue
	if err := json.Unmarshal(in.Body.Value, &fv); err != nil {
		return nil, huma.Error400BadRequest("invalid value JSON")
	}
	if !fv.Validate() {
		return nil, huma.Error422UnprocessableEntity("value must be exactly one of bool/number/string")
	}

	if err := h.svc.SetValue(ctx, in.Body.Key, fv); err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &SetVarResp{Body: SetVarRespBody{Message: "var successfully updated"}}, nil
}
