package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	grpc_handler "github.com/sachinggsingh/e-comm/internal/api/grpc"
	"github.com/sachinggsingh/e-comm/internal/api/restapi"
	cache "github.com/sachinggsingh/e-comm/internal/caches"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/middleware"
	"github.com/sachinggsingh/e-comm/internal/service"
	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
)

type Server struct {
	env        *config.Env
	db         *db.Database
	redisCache *cache.RedisCache
}

// NewServer accepts an externally created RedisCache so the client is
// created/closed by the caller (avoids duplicate connections).
func NewServer(env *config.Env, database *db.Database, redisCache *cache.RedisCache) *Server {
	return &Server{
		env:        env,
		db:         database,
		redisCache: redisCache,
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
	// Create auth middleware
	authMiddleware := &middleware.AuthMiddleware{Redis: s.redisCache}

	// Public routes (no authentication required)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		userHandler := restapi.NewUserHandler(userService, s.redisCache, s.env)
		userHandler.Register(w, r)
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		userHandler := restapi.NewUserHandler(userService, s.redisCache, s.env)
		userHandler.Login(w, r)
	})

	// Protected routes (authentication required)
	// Profile endpoint with caching middleware
	http.Handle("/profile", authMiddleware.Validate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userHandler := restapi.NewUserHandler(userService, s.redisCache, s.env)
		userHandler.Profile(w, r)
	})))
}
