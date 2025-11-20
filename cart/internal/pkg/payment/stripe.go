package payment

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

// PaymentItem represents a product item in the payment
type PaymentItem struct {
	Name        string
	Description string
	Price       float64
	Quantity    int64
}

type PaymentClient interface {
	CreatePayment(items []PaymentItem, userID string, orderID string) (*stripe.CheckoutSession, error)
	CheckPaymentStatus(pID string) (*stripe.CheckoutSessionStatus, error)
}

type Payment struct {
	stripeSecretKey string
	successUrl      string
	failureUrl      string
}

func NewPaymentClient(stripeSecretKey, successUrl, failureUrl string) PaymentClient {
	return &Payment{
		stripeSecretKey: stripeSecretKey,
		successUrl:      successUrl,
		failureUrl:      failureUrl,
	}
}

func (p *Payment) CreatePayment(items []PaymentItem, userID string, orderID string) (*stripe.CheckoutSession, error) {
	if len(items) == 0 {
		return nil, errors.New("payment items cannot be empty")
	}

	stripe.Key = p.stripeSecretKey

	// Build line items from product details
	lineItems := make([]*stripe.CheckoutSessionLineItemParams, 0, len(items))
	for _, item := range items {
		if item.Price <= 0 {
			return nil, fmt.Errorf("invalid price for item %s: %f", item.Name, item.Price)
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for item %s: %d", item.Name, item.Quantity)
		}

		amountInCents := int64(item.Price * 100)

		productData := &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
			Name: stripe.String(item.Name),
		}

		// Add description if provided
		if item.Description != "" {
			productData.Description = stripe.String(item.Description)
		}

		lineItem := &stripe.CheckoutSessionLineItemParams{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				UnitAmount:  stripe.Int64(amountInCents),
				Currency:    stripe.String("usd"),
				ProductData: productData,
			},
			Quantity: stripe.Int64(item.Quantity),
		}
		lineItems = append(lineItems, lineItem)
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
		Mode:               stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:         stripe.String(p.successUrl),
		CancelURL:          stripe.String(p.failureUrl),
	}
	params.AddMetadata("order_id", orderID)
	params.AddMetadata("user_id", userID)

	session, err := session.New(params)
	if err != nil {
		return nil, fmt.Errorf("payment session creation failed: %w", err)
	}
	return session, nil
}
func (p *Payment) CheckPaymentStatus(pID string) (*stripe.CheckoutSessionStatus, error) {
	stripe.Key = p.stripeSecretKey

	session, err := session.Get(pID, nil)
	if err != nil {
		return nil, err
	}

	status := session.Status
	fmt.Println(&status, "  Checking the & status")
	return &status, nil
}
