package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/sachinggsingh/e-comm/internal/helper"
	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var newUser model.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := u.userService.RegisterUser(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input data: "+err.Error(), http.StatusBadRequest)
		return
	}
	loggedInUser, err := u.userService.LoginUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := model.LoginResponse{
		ID:           loggedInUser.User_id,
		Email:        *loggedInUser.Email,
		Token:        *loggedInUser.Token,
		RefreshToken: *loggedInUser.Refresh_Token,
		User_id:      loggedInUser.User_id,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userIDFromRequest := r.URL.Query().Get("user_id")
	if userIDFromRequest == "" {
		http.Error(w, "User id is required", http.StatusBadRequest)
		return
	}
	userID, err := helper.Authorize(r, userIDFromRequest)
	if err != nil {
		http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
		return
	}
	user := &model.User{User_id: userID}
	result, err := u.userService.Profile(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch profile: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)

}

// profile is not working till now perfectly
