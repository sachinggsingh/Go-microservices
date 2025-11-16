package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sachinggsingh/e-comm/internal/errors"
	"github.com/sachinggsingh/e-comm/internal/model"
	"github.com/sachinggsingh/e-comm/internal/pkg"
	"github.com/sachinggsingh/e-comm/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartService struct {
	cartRepo      repository.CartRepository
	productClient *pkg.ProductClient
}

func NewCartService(cartRepo repository.CartRepository, productClient *pkg.ProductClient) *CartService {
	return &CartService{
		cartRepo:      cartRepo,
		productClient: productClient,
	}
}

// CalculateTotal calculates the total price for a cart item
func CalculateTotal(price float64, quantity int) float64 {
	return price * float64(quantity)
}

// CalculateCartTotal calculates the total amount for the entire cart
func CalculateCartTotal(items []model.CartItem) float64 {
	var total float64
	for _, item := range items {
		total += item.Total
	}
	return total
}

// ValidateCartItems validates cart items
func validateCartItems(items []model.CartItem) error {
	if len(items) == 0 {
		return errors.ErrEmptyCart
	}

	for _, item := range items {
		if item.Product_id == "" {
			return errors.ErrInvalidItem
		}
		if item.Quantity <= 0 {
			return errors.ErrInvalidQuantity
		}
		if item.Price <= 0 {
			return errors.ErrInvalidPrice
		}
	}

	return nil
}

// ValidateProductsWithGRPC validates products exist and optionally updates prices from product service
func (c *CartService) ValidateProductsWithGRPC(ctx context.Context, items []model.CartItem) ([]model.CartItem, error) {
	if c.productClient == nil {
		// If product client is not available, skip validation
		log.Println("Product client not available, skipping product validation")
		return items, nil
	}

	validatedItems := make([]model.CartItem, 0, len(items))
	for _, item := range items {
		// Validate product exists via gRPC
		product, err := c.productClient.GetProduct(ctx, item.Product_id)
		if err != nil {
			return nil, fmt.Errorf("product %s validation failed: %w", item.Product_id, err)
		}

		// Use price from product service to ensure consistency
		validatedItem := item
		validatedItem.Price = product.Price
		validatedItem.Total = CalculateTotal(product.Price, item.Quantity)
		validatedItems = append(validatedItems, validatedItem)
	}

	return validatedItems, nil
}

func (c *CartService) CreateCart(userID string, items []model.CartItem) (*model.Cart, error) {
	if userID == "" {
		return nil, errors.ErrInvalidUserID
	}

	// Validate items
	if err := validateCartItems(items); err != nil {
		return nil, err
	}

	// Validate products via gRPC and get updated prices
	ctx := context.Background()
	validatedItems, err := c.ValidateProductsWithGRPC(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	// Check if cart already exists for user
	existingCart, err := c.cartRepo.FindCartByUserID(userID)
	if err != nil && err != errors.ErrCartNotFound {
		log.Printf("Error checking existing cart: %v", err)
		return nil, err
	}

	now := time.Now().UTC()

	// If cart exists, add items to existing cart
	if existingCart != nil {
		return c.AddItemsToCart(userID, validatedItems)
	}

	// Create new cart
	cartID := primitive.NewObjectID()
	cartItems := make([]model.CartItem, 0, len(validatedItems))

	for _, item := range validatedItems {
		cartItemID := primitive.NewObjectID()
		cartItem := model.CartItem{
			ID:          cartItemID,
			Product_id:  item.Product_id,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Total:       item.Total,
			CartItem_id: cartItemID.Hex(),
		}
		cartItems = append(cartItems, cartItem)
	}

	totalAmount := CalculateCartTotal(cartItems)

	cart := &model.Cart{
		ID:          cartID,
		User_id:     userID,
		Items:       cartItems,
		TotalAmount: totalAmount,
		Created_at:  now,
		Updated_at:  now,
		Cart_id:     cartID.Hex(),
	}

	return c.cartRepo.CreateCart(cart)
}

func (c *CartService) FindCartByUserID(userID string) (*model.Cart, error) {
	if userID == "" {
		return nil, errors.ErrInvalidUserID
	}
	return c.cartRepo.FindCartByUserID(userID)
}

func (c *CartService) UpdateCart(userID string, items []model.CartItem) (*model.Cart, error) {
	if userID == "" {
		return nil, errors.ErrInvalidUserID
	}

	// Validate items
	if err := validateCartItems(items); err != nil {
		return nil, err
	}

	// Validate products via gRPC and get updated prices
	ctx := context.Background()
	validatedItems, err := c.ValidateProductsWithGRPC(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	// Get existing cart
	existingCart, err := c.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update items
	cartItems := make([]model.CartItem, 0, len(validatedItems))
	for _, item := range validatedItems {
		cartItemID := primitive.NewObjectID()
		cartItem := model.CartItem{
			ID:          cartItemID,
			Product_id:  item.Product_id,
			Price:       item.Price,
			Quantity:    item.Quantity,
			Total:       item.Total,
			CartItem_id: cartItemID.Hex(),
		}
		cartItems = append(cartItems, cartItem)
	}

	totalAmount := CalculateCartTotal(cartItems)

	updatedCart := &model.Cart{
		ID:          existingCart.ID,
		User_id:     userID,
		Items:       cartItems,
		TotalAmount: totalAmount,
		Created_at:  existingCart.Created_at,
		Updated_at:  time.Now().UTC(),
		Cart_id:     existingCart.Cart_id,
	}

	return c.cartRepo.UpdateCart(updatedCart)
}

func (c *CartService) DeleteCart(userID string) error {
	if userID == "" {
		return errors.ErrInvalidUserID
	}
	return c.cartRepo.DeleteCart(userID)
}

func (c *CartService) AddItemsToCart(userID string, items []model.CartItem) (*model.Cart, error) {
	if userID == "" {
		return nil, errors.ErrInvalidUserID
	}

	if err := validateCartItems(items); err != nil {
		return nil, err
	}

	// Validate products via gRPC and get updated prices
	ctx := context.Background()
	validatedItems, err := c.ValidateProductsWithGRPC(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("product validation failed: %w", err)
	}

	existingCart, err := c.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Add new items to existing items
	for _, item := range validatedItems {
		// 	Check if product already exists in cart
		found := false
		for i, existingItems := range existingCart.Items {
			if existingItems.Product_id == item.Product_id {
				existingCart.Items[i].Quantity += item.Quantity
				existingCart.Items[i].Total = CalculateTotal(existingCart.Items[i].Price, existingCart.Items[i].Quantity)
				found = true
				break
			}
		}
		if !found {
			cartItemID := primitive.NewObjectID()
			total := CalculateTotal(item.Price, item.Quantity)
			cartItem := model.CartItem{
				ID:          cartItemID,
				Product_id:  item.Product_id,
				Price:       item.Price,
				Quantity:    item.Quantity,
				Total:       total,
				CartItem_id: cartItemID.Hex(),
			}
			existingCart.Items = append(existingCart.Items, cartItem)
		}
	}

	// Recalculate total
	existingCart.TotalAmount = CalculateCartTotal(existingCart.Items)
	existingCart.Updated_at = time.Now().UTC()

	return c.cartRepo.UpdateCart(existingCart)
}

func (c *CartService) UpdateCartItem(userID string, productID string, quantity int) (*model.Cart, error) {
	if userID == "" {
		return nil, errors.ErrInvalidUserID
	}
	if productID == "" {
		return nil, errors.ErrInvalidItem
	}
	if quantity <= 0 {
		return nil, errors.ErrInvalidQuantity
	}

	// Get existing cart
	existingCart, err := c.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Find and update item
	found := false
	for i, item := range existingCart.Items {
		if item.Product_id == productID {
			existingCart.Items[i].Quantity = quantity
			existingCart.Items[i].Total = CalculateTotal(item.Price, quantity)
			found = true
			break
		}
	}

	if !found {
		return nil, errors.ErrInvalidItem
	}

	// Recalculate total
	existingCart.TotalAmount = CalculateCartTotal(existingCart.Items)
	existingCart.Updated_at = time.Now().UTC()

	return c.cartRepo.UpdateCart(existingCart)
}

func (c *CartService) RemoveItemFromCart(userID string, productID string) (*model.Cart, error) {
	if userID == "" {
		return nil, errors.ErrInvalidUserID
	}
	if productID == "" {
		return nil, errors.ErrInvalidItem
	}

	// Get existing cart
	existingCart, err := c.cartRepo.FindCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Remove item
	newItems := make([]model.CartItem, 0)
	found := false
	for _, item := range existingCart.Items {
		if item.Product_id != productID {
			newItems = append(newItems, item)
		} else {
			found = true
		}
	}

	if !found {
		return nil, errors.ErrInvalidItem
	}

	existingCart.Items = newItems
	existingCart.TotalAmount = CalculateCartTotal(existingCart.Items)
	existingCart.Updated_at = time.Now().UTC()

	return c.cartRepo.UpdateCart(existingCart)
}

func (c *CartService) ClearCart(userID string) error {
	if userID == "" {
		return errors.ErrInvalidUserID
	}
	return c.cartRepo.DeleteCart(userID)
}
