package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sachinggsingh/e-comm/internal/api/restapi"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/helper"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/service"
)

type Server struct {
	env *config.Env
	db  *db.Database
}

func NewServer(env *config.Env, database *db.Database) *Server {
	return &Server{
		env: env,
		db:  database,
	}
}

func (s *Server) StartServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Welcome to the UserRoute"))
	})

	addr := fmt.Sprintf(":%s", s.env.PORT)
	log.Printf("Starting server on port %s\n", s.env.PORT)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server")
	}
}

func (s *Server) UserRoutes(userService *service.UserService) {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		userHandler := restapi.NewUserHandler(userService)
		userHandler.Register(w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler := restapi.NewUserHandler(userService)
		userHandler.Login(w, r)
	})
	// http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
	// 	userHandler := restapi.NewUserHandler(userService)
	// 	userHandler.Logout(w, r)
	// })

	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		_, err := helper.Authorize(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		userHanlder := restapi.NewUserHandler(userService)
		userHanlder.Profile(w, r)
	})
}
