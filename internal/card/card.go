package card

import (
	"github.com/stripe/stripe-go/paymentintent"
	"github.com/stripe/stripe-go/v72"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Key

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount: stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// NOTE: you can add info to this transaction. See below
	// params.AddMetadata("key", "value")

	pi, err := paymentintent.New(params)
	if err != nil {
		var msg = ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = getCardErrorMessage(stripeErr)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

func getCardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card has expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC code"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect zip/postal code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The amount is too large to charge from your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The amount is too small to charge from your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your postal code is invalid"
	default:
		msg = "Your card was declined"
	}

	return msg
}
