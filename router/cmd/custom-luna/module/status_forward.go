package module

import (
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/wundergraph/cosmo/router/core"
)

func init() {
	core.RegisterModule(&StatusForwardModule{})
}

const StatusForwardModuleID = "com.getluna.status-forward"

// forwardedStatusCodes lists the GraphQL error extension statusCodes that should
// override the HTTP response status. Subgraphs return these but the upstream
// router does not propagate them by default.
var forwardedStatusCodes = map[int64]bool{
	498: true,
	409: true,
}

type StatusForwardModule struct{}

func (m *StatusForwardModule) Middleware(ctx core.RequestContext, next http.Handler) {
	sw := &statusForwardWriter{ResponseWriter: ctx.ResponseWriter()}
	next.ServeHTTP(sw, ctx.Request())
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
// code in errors[].extensions.statusCode and applies it before delegating.
type statusForwardWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (s *statusForwardWriter) WriteHeader(code int) {
	if s.wroteHeader {
		return
	}
	s.wroteHeader = true
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusForwardWriter) Write(p []byte) (int, error) {
	if !s.wroteHeader {
		if code := extractForwardedStatus(p); code != 0 {
			s.ResponseWriter.WriteHeader(int(code))
			s.wroteHeader = true
		}
	}
	return s.ResponseWriter.Write(p)
}

func extractForwardedStatus(body []byte) int64 {
	var found int64
	_, _ = jsonparser.ArrayEach(body, func(value []byte, _ jsonparser.ValueType, _ int, _ error) {
		if found != 0 {
			return
		}
		code, err := jsonparser.GetInt(value, "extensions", "statusCode")
		if err == nil && forwardedStatusCodes[code] {
			found = code
		}
	}, "errors")
	return found
}

var _ core.RouterMiddlewareHandler = (*StatusForwardModule)(nil)
