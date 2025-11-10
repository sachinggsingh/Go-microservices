package repository

import (
	"context"
	"time"

	"github.com/sachinggsingh/e-comm/internal/intra/db"
	"github.com/sachinggsingh/e-comm/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Register(user *model.User) error
	Login(user *model.User) (*model.User, error)
	// Logout(user *model.User) error
	Profile(user *model.User) (*model.User, error)
	CheckIfEmailExist(email string) (bool, error)
	FindUserByEmail(email string) (*model.User, error)
}

type userRepo struct {
	userColl *mongo.Collection
}

func NewUserRepository(database *db.Database) UserRepository {
	return &userRepo{
		userColl: database.UserCollection,
	}
}

func (u *userRepo) Register(user *model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := u.userColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *userRepo) Login(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": user.User_id,
	}
	update := bson.M{
		"$set": bson.M{
			"token":         user.Token,
			"refresh_token": user.Refresh_Token,
		},
	}
	result := u.userColl.FindOneAndUpdate(ctx, filter, update)
	if result.Err() != nil {
		return nil, result.Err()
	}
	return user, nil
}

//	func (u *userRepo) Logout(user *model.User) error {
//		return nil
//	}
func (u *userRepo) Profile(user *model.User) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": user.User_id,
	}
	var result model.User
	err := u.userColl.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &result, nil
}

func (u *userRepo) CheckIfEmailExist(email string) (bool, error) {
	filter := bson.M{
		"email": email,
	}
	var result model.User
	err := u.userColl.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func (u *userRepo) FindUserByEmail(email string) (*model.User, error) {
	filter := bson.M{
		"email": email,
	}
	var result model.User
	err := u.userColl.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, nil
	}
	return &result, nil
}
