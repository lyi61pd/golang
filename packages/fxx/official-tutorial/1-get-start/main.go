package main

import (
	"context"
	"go.uber.org/fx"
	"io"
	"net"
	"fmt"
	"net/http"
	"go.uber.org/zap"
)

func NewHTTPServer(lc fx.Lifecycle, mux *http.ServeMux, log *zap.Logger) *http.Server {
   srv := &http.Server{Addr: ":8080", Handler: mux}
   lc.Append(fx.Hook{
       OnStart: func(ctx context.Context) error {
           ln, err := net.Listen("tcp", srv.Addr)
           if err != nil {
               return err
           }
           log.Info("Starting HTTP server", zap.String("addr", srv.Addr))
           go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   if _, err := io.Copy(w, r.Body); err != nil {
       h.log.Warn("Failed to handle request", zap.Error(err))
   }
}

func (*EchoHandler) Pattern() string {
   return "/echo"
}

// Route is an http.Handler that knows the mux pattern
// under which it will be registered.
type Route interface {
   http.Handler

   // Pattern reports the path at which this is registered.
   Pattern() string
}

func NewServeMux(routes []Route) *http.ServeMux {
   mux := http.NewServeMux()
   for _, route := range routes {
       mux.Handle(route.Pattern(), route)
   }
   return mux
}

type EchoHandler struct {
   log *zap.Logger
}

func NewEchoHandler(log *zap.Logger) *EchoHandler {
   return &EchoHandler{log: log}
}

// HelloHandler is an HTTP handler that
// prints a greeting to the user.
type HelloHandler struct {
   log *zap.Logger
}

// NewHelloHandler builds a new HelloHandler.
func NewHelloHandler(log *zap.Logger) *HelloHandler {
   return &HelloHandler{log: log}
}

func (*HelloHandler) Pattern() string {
   return "/hello"
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
   body, err := io.ReadAll(r.Body)
   if err != nil {
       h.log.Error("Failed to read request", zap.Error(err))
       http.Error(w, "Internal server error", http.StatusInternalServerError)
       return
   }

   if _, err := fmt.Fprintf(w, "Hello, %s\n", body); err != nil {
       h.log.Error("Failed to write response", zap.Error(err))
       http.Error(w, "Internal server error", http.StatusInternalServerError)
       return
   }
}

// AsRoute annotates the given constructor to state that
// it provides a route to the "routes" group.
func AsRoute(f any) any {
   return fx.Annotate(
       f,
       fx.As(new(Route)),
       fx.ResultTags(`group:"routes"`),
   )
}

func main() {
	fx.New(
		fx.Provide(
			NewHTTPServer,
			AsRoute(NewEchoHandler),
           	AsRoute(NewHelloHandler),
			fx.Annotate(
				NewServeMux,
				fx.ParamTags(`group:"routes"`),
			),
			zap.NewExample,
		),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
