package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/sachinggsingh/e-comm/internal/helper"
	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (u *UserService) RegisterUser(user *model.User) error {
	if user.Email == nil || user.Password == nil {
		return errors.New("email and password are required")
	}

	exist, err := u.repo.CheckIfEmailExist(*user.Email)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if exist {
		return errors.New("user already exists")
	}

	hashedPassword, err := helper.HashPassword(*user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	now := time.Now()
	user.ID = primitive.NewObjectID()
	user.User_id = user.ID.Hex()

	token, refreshToken, err := helper.GenerateToken(user.User_id, *user.Email)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	newUser := &model.User{
		ID:            user.ID,
		Password:      &hashedPassword,
		Email:         user.Email,
		Token:         &token,
		Refresh_Token: &refreshToken,
		Created_at:    now,
		Updated_at:    now,
		User_id:       user.User_id,
	}
	return u.repo.Register(newUser)
}

func (u *UserService) LoginUser(user *model.User) (*model.User, error) {
	//checking the user Exits
	if user.Email == nil || user.Password == nil {
		return nil, errors.New("email and password are required")
	}
	storedUser, err := u.repo.FindUserByEmail(*user.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	if storedUser == nil {
		return nil, errors.New("user not found")
	}
	isPasswordOk := helper.CheckPasswordHash(*user.Password, *storedUser.Password)
	if !isPasswordOk {
		return nil, errors.New("password is incorrect")
	}
	token, refreshToken, err := helper.GenerateToken(storedUser.User_id, *storedUser.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	storedUser.Token = &token
	storedUser.Refresh_Token = &refreshToken
	return u.repo.Login(storedUser)
}

// func (u *UserService) LogoutUser(user *model.User) error {
// 	return u.repo.Logout(user)
// }

func (u *UserService) Profile(user *model.User) (*model.User, error) {
	return u.repo.Profile(user)
}
