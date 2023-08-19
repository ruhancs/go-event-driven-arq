package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router chi.Router
	Handlers map[string]http.HandlerFunc //map de string retorna http.handlefunc
	WebServerPort string 
}

func NewWebServer (webServerPort string) *WebServer {
	return &WebServer{
		Router: chi.NewRouter(),
		Handlers: make(map[string]http.HandlerFunc),
		WebServerPort: webServerPort,
	}
}

func (server *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	//adiciona um handler com o nome dele que sera o path com a funcao que ele executara que é o handler
	server.Handlers[path] = handler 
}

func (server *WebServer) Start() {
	//utilizando um middleware de logs
	server.Router.Use(middleware.Logger)
	//cada handler inserido sera registrado no chi
	for path, handler := range server.Handlers {
		server.Router.Post(path,handler)
	}

	http.ListenAndServe(server.WebServerPort,server.Router)
}