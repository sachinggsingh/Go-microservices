package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	grpc_handler "github.com/sachinggsingh/e-comm/internal/api/grpc"
	"github.com/sachinggsingh/e-comm/internal/api/restapi"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/helper"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/service"
	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
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

func (s *Server) StartServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello Welcome to the UserRoute"))
	})

	addr := fmt.Sprintf(":%s", s.env.PORT)
	log.Printf("Starting HTTP server on port %s\n", s.env.PORT)

	if err := http.ListenAndServe(addr, nil); err != nil {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

func StartGRPC() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf(" gRPC failed to listen on :9090: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterValidateTokenServer(grpcServer, &grpc_handler.AuthServer{})

	log.Println("=" + strings.Repeat("=", 50) + "=")
	log.Println("Auth gRPC Server is running on :9090")
	// log.Println("   Ready to accept ValidateToken requests")
	log.Println("=" + strings.Repeat("=", 50) + "=")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" gRPC server failed to serve: %v", err)
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
