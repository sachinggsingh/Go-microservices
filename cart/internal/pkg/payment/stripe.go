package payment

import (
	"errors"
	"fmt"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

type PaymentClient interface {
	CreatePayment(amount float64, userID string, orderID string) (*stripe.CheckoutSession, error)
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

func (p *Payment) CreatePayment(amount float64, userID string, orderID string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := int64(amount * 100)

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					UnitAmount: stripe.Int64(amountInCents),
					Currency:   stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Cloth"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(p.successUrl),
		CancelURL:  stripe.String(p.failureUrl),
	}
	params.AddMetadata("order_id", orderID)
	params.AddMetadata("user_id", userID)

	session, err := session.New(params)
	if err != nil {
		return nil, errors.New("Payment Creations session failed")
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
