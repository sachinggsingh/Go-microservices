package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sachinggsingh/e-comm/internal/api/restapi"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/service"
)

type Server struct {
	env *config.Env
	db  *db.Database
	r   *mux.Router
}

func NewServer(env *config.Env, db *db.Database) *Server {
	return &Server{
		env: env,
		db:  db,
		r:   mux.NewRouter(),
	}
}

func (s *Server) StartServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Welcome to the CartRoute"))
	})

	port := fmt.Sprintf(":%s", s.env.PORT)
	log.Printf("Starting server on port %s\n", s.env.PORT)

	if err := http.ListenAndServe(port, s.r); err != nil {
		log.Fatalf("Failed to start server")
	}
}

func (s *Server) CartRoute(cartservice *service.CartService) {
	cartHandler := restapi.NewCartHandler(cartservice)

	s.r.HandleFunc("/cart", cartHandler.CreateCart).Methods("POST")
	s.r.HandleFunc("/cart", cartHandler.GetAllCarts).Methods("GET")
	s.r.HandleFunc("/cart/{cart_id}", cartHandler.GetCartById).Methods("GET")
}
