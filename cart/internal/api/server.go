package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sachinggsingh/e-comm/internal/api/restapi"
	"github.com/sachinggsingh/e-comm/internal/config"
	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/middleware"
	"github.com/sachinggsingh/e-comm/internal/pkg/payment"
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
	port := fmt.Sprintf(":%s", s.env.PORT)
	log.Printf("Starting server on port %s\n", s.env.PORT)

	if err := http.ListenAndServe(port, s.r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func (s *Server) CartRoute(cartservice *service.CartService, paymentClient payment.PaymentClient) {
	cartHandler := restapi.NewCartHandler(cartservice, paymentClient)

	// Create cart - requires authentication
	s.r.Handle("/cart", middleware.GetUserIdFromToken(http.HandlerFunc(cartHandler.CreateCart))).Methods("POST")

	// Get cart by user ID - requires authentication
	s.r.Handle("/cart/{user_id}", middleware.GetUserIdFromToken(http.HandlerFunc(cartHandler.FindCartByUserID))).Methods("GET")

	// Update cart - requires authentication
	s.r.Handle("/cart/{user_id}", middleware.GetUserIdFromToken(http.HandlerFunc(cartHandler.UpdateCart))).Methods("PUT")

	// Delete cart - requires authentication
	s.r.Handle("/cart/{user_id}", middleware.GetUserIdFromToken(http.HandlerFunc(cartHandler.DeleteCart))).Methods("DELETE")

	// Checkout cart - creates Stripe payment session - requires authentication
	s.r.Handle("/cart/checkout", middleware.GetUserIdFromToken(http.HandlerFunc(cartHandler.CheckoutCart))).Methods("POST")
}
