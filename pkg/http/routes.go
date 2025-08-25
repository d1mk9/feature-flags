package http

import (
	"feature-flags/pkg/handlers"
	"feature-flags/pkg/service"

	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API, svc service.Flags) {
	h := handlers.NewFeatureHandler(svc)

	huma.Get(api, "/var/{var_name}", h.GetVar,
		func(op *huma.Operation) { op.Description = "Get variable value (15 min in-memory cache)" },
	)
	huma.Post(api, "/var/set", h.SetVar,
		func(op *huma.Operation) { op.Description = "Set variable value and invalidate cache" },
	)
}
