package module

import (
	"context"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/wundergraph/cosmo/router/core"
	"go.opentelemetry.io/otel/attribute"
)

func init() {
	core.RegisterModule(&StatusForwardModule{})
}

const (
	StatusForwardModuleID    = "com.getluna.status-forward"
	statusForwardTracerName  = "lunacare/cosmo/router/status-forward"
	statusForwardSpanName    = "lunacare.status_forward"
	statusForwardAppliedAttr = "lunacare.status_forward.applied"
	statusForwardCodeAttr    = "lunacare.status_forward.code"
)

// defaultForwardedCodes are the GraphQL error extension statusCodes that should
// override the HTTP response status by default. Subgraphs return these but the
// upstream router does not propagate them. Override via config:
//
//	modules:
//	  com.getluna.status-forward:
//	    forwarded_codes: [498, 409, 429]
var defaultForwardedCodes = []int64{498, 409}

type StatusForwardModule struct {
	// Enabled toggles the module. Defaults to false — status forwarding is a
	// no-op unless the router config explicitly opts in:
	//
	//	modules:
	//	  com.getluna.status-forward:
	//	    enabled: true
	Enabled bool `mapstructure:"enabled"`

	// ForwardedCodes overrides the list of statusCode values lifted from
	// errors[].extensions.statusCode into the HTTP response status. Empty ⇒
	// defaultForwardedCodes ([498, 409]).
	ForwardedCodes []int64 `mapstructure:"forwarded_codes"`

	forwardSet map[int64]bool
}

func (m *StatusForwardModule) Provision(_ *core.ModuleContext) error {
	codes := m.ForwardedCodes
	if len(codes) == 0 {
		codes = defaultForwardedCodes
	}
	m.forwardSet = make(map[int64]bool, len(codes))
	for _, c := range codes {
		m.forwardSet[c] = true
	}
	return nil
}

func (m *StatusForwardModule) Middleware(ctx core.RequestContext, next http.Handler) {
	r := ctx.Request()
	w := ctx.ResponseWriter()

	if !m.Enabled {
		next.ServeHTTP(w, r)
		return
	}

	sw := &statusForwardWriter{
		ResponseWriter: w,
		codes:          m.forwardSet,
		reqCtx:         r.Context(),
	}
	next.ServeHTTP(sw, r)
}

func (m *StatusForwardModule) Module() core.ModuleInfo {
	return core.ModuleInfo{
		ID: StatusForwardModuleID,
		New: func() core.Module {
			return &StatusForwardModule{}
		},
	}
}

// statusForwardWriter intercepts the first Write to look for a forwarded status
// code in errors[].extensions.statusCode and applies it before delegating. The
// inspection emits a short-lived span the first time bytes flow through, so
// every request gets one statusForwardSpanName span with an `applied` attribute
// even when no forwarding fires.
type statusForwardWriter struct {
	http.ResponseWriter
	codes       map[int64]bool
	reqCtx      context.Context
	wroteHeader bool
	scanned     bool
}

func (s *statusForwardWriter) WriteHeader(code int) {
	if s.wroteHeader {
		return
	}
	s.wroteHeader = true
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusForwardWriter) Write(p []byte) (int, error) {
	if !s.scanned {
		s.scanAndForward(p)
		s.scanned = true
	}
	return s.ResponseWriter.Write(p)
}

func (s *statusForwardWriter) scanAndForward(p []byte) {
	tracer := tracerFromContext(s.reqCtx, statusForwardTracerName)
	_, span := tracer.Start(s.reqCtx, statusForwardSpanName)
	defer span.End()
	span.SetAttributes(attribute.String("module.id", StatusForwardModuleID))

	code := extractForwardedStatus(p, s.codes)
	applied := code != 0 && !s.wroteHeader
	span.SetAttributes(attribute.Bool(statusForwardAppliedAttr, applied))
	if applied {
		span.SetAttributes(attribute.Int64(statusForwardCodeAttr, code))
		s.ResponseWriter.WriteHeader(int(code))
		s.wroteHeader = true
	}
}

func extractForwardedStatus(body []byte, codes map[int64]bool) int64 {
	var found int64
	_, _ = jsonparser.ArrayEach(body, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		if found != 0 {
			return
		}
		code, err := jsonparser.GetInt(value, "extensions", "statusCode")
		if err == nil && codes[code] {
			found = code
		}
	}, "errors")
	return found
}

var (
	_ core.RouterMiddlewareHandler = (*StatusForwardModule)(nil)
	_ core.Provisioner             = (*StatusForwardModule)(nil)
)
