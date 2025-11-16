package api

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	grpc_handler "github.com/sachinggsingh/e-comm/internal/api/grpc"
	"github.com/sachinggsingh/e-comm/internal/api/restapi"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/service"
	proto "github.com/sachinggsingh/e-comm/pb"
	"google.golang.org/grpc"
)

type Server struct {
	env *config.Env
	db  *db.Database
	r   *mux.Router
}

func NewServer(env *config.Env, database *db.Database) *Server {
	return &Server{
		env: env,
		db:  database,
		r:   mux.NewRouter(),
	}
}

func (s *Server) StartServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello welcome to product route"))
	})

	addr := fmt.Sprintf(":%s", s.env.PORT)
	log.Printf("Starting server on port %s\n", s.env.PORT)

	if err := http.ListenAndServe(addr, s.r); err != nil {
		log.Fatalf("Failed to start server")
	}
	return nil
}

func (s *Server) GrpcServer() error {
	lis, err := net.Listen("tcp", ":9091")
	if err != nil {
		log.Fatalf(" gRPC failed to listen on :9091: %v", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterGetProductsServer(grpcServer, grpc_handler.NewProductServer(s.db))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(" gRPC server failed to serve: %v", err)
	}
	log.Printf("GRPC server running")

	return nil

}

func (s *Server) ProductRoutes(productService *service.Productservice) {
	productHandler := restapi.NewProductHandler(productService)

	s.r.HandleFunc("/product", productHandler.CreateProduct).Methods("POST")
	s.r.HandleFunc("/product", productHandler.GetAllProducts).Methods("GET")
	s.r.HandleFunc("/product/{product_id}", productHandler.GetProductById).Methods("GET")
	// http.HandleFunc("/updateproduct", func(w http.ResponseWriter, r *http.Request) {
	// 	productHanlder := restapi.NewProductHandler(productService)
	// 	productHanlder.UpdateProduct(w, r)
	// })
	// http.HandleFunc("/deleteproduct", func(w http.ResponseWriter, r *http.Request) {
	// 	productHanlder := restapi.NewProductHandler(productService)
	// 	productHanlder.DeleteProduct(w, r)
	// })
}
