package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"feature-flags/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

type VarsHandler struct {
	Svc service.Vars
}

// ==== I/O схемы ====

type GetVarInput struct {
	VarName string `path:"var_name"`
}
type GetVarResp struct {
	Body struct {
		Key   string      `json:"key"`
		Value interface{} `json:"value"`
	}
}

type SetVarReq struct {
	Body struct {
		Key   string          `json:"key"   required:"true"`
		Value json.RawMessage `json:"value" required:"true"`
	}
}
type SetVarResp struct {
	Body struct {
		Message string `json:"message"`
	}
}

// ==== Handlers ====

func (h *VarsHandler) GetVar(ctx context.Context, in *GetVarInput) (*GetVarResp, error) {
	raw, err := h.Svc.GetValue(in.VarName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, huma.Error404NotFound("variable not found")
		}
		return nil, huma.Error500InternalServerError(fmt.Sprintf("db: %v", err))
	}

	var any interface{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &any)
	}

	out := &GetVarResp{}
	out.Body.Key = in.VarName
	out.Body.Value = any
	return out, nil
}

func (h *VarsHandler) SetVar(ctx context.Context, in *SetVarReq) (*SetVarResp, error) {
	if in == nil || in.Body.Key == "" {
		return nil, huma.Error400BadRequest("key is required")
	}
	if err := h.Svc.SetValue(in.Body.Key, in.Body.Value); err != nil {
		return nil, huma.Error500InternalServerError(fmt.Sprintf("db: %v", err))
	}

	out := &SetVarResp{}
	out.Body.Message = "var successfully updated"
	return out, nil
}
