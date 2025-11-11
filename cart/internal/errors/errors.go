package errors

import "errors"

var (
	ErrCartNotFound    = errors.New("cart not found")
	ErrInvalidItem     = errors.New("invalid cart item")
	ErrInvalidQuantity = errors.New("quantity must be greater than 0")
	ErrInvalidPrice    = errors.New("price must be greater than 0")
	ErrEmptyCart       = errors.New("cart items cannot be empty")
	ErrInvalidUserID   = errors.New("user_id is required")
	ErrNoTokenProvided = errors.New("no token provided")
	ErrInvalidToken    = errors.New("invalid token")
	ErrInGenerating    = errors.New("error in generating token")
	DifferentTokenUsed = errors.New("uid missing: likely refresh token used instead of access token")
)
