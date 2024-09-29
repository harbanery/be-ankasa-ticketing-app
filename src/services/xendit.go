package services

import (
	"ankasa-be/src/models"
	"context"
	"fmt"
	"os"

	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/payment_method"
)

var Client *xendit.APIClient

func InitXendit() {
	secretKey := os.Getenv("XENDIT_SECRET_KEY")

	Client = xendit.NewClient(secretKey)
}

func EWalletPaymentMethod(ewallet *models.PaymentMethodEWallet) (*string, error) {
	paymentMethodParameters := *payment_method.
		NewPaymentMethodParameters(
			payment_method.PaymentMethodType("EWALLET"),
			payment_method.PaymentMethodReusability("ONE_TIME_USE"))

	successReturnUrl := "https://redirect.me/goodstuff"
	failureReturnUrl := "https://redirect.me/badstuff"
	cancelReturnUrl := "https://redirect.me/nostuff"
	pendingReturnUrl := "https://redirect.me/nostuff"

	if valid := payment_method.EWalletChannelCode(ewallet.Name).IsValid(); !valid {
		return nil, fmt.Errorf("your payment method not valid")
	}

	ewalletParams := payment_method.EWalletParameters{
		ChannelCode: payment_method.EWalletChannelCode(ewallet.Name),
		ChannelProperties: &payment_method.EWalletChannelProperties{
			SuccessReturnUrl: &successReturnUrl,
			FailureReturnUrl: &failureReturnUrl,
			CancelReturnUrl:  &cancelReturnUrl,
			PendingReturnUrl: &pendingReturnUrl,
			MobileNumber:     ewallet.MobileNumber,
			Cashtag:          ewallet.Cashtag,
		},
	}

	paymentMethodParameters.SetEwallet(ewalletParams)

	if _, ok := paymentMethodParameters.GetEwalletOk(); !ok {
		return nil, fmt.Errorf("ewallet param error")
	}

	createPMResp, _, err := Client.PaymentMethodApi.
		CreatePaymentMethod(context.Background()).
		PaymentMethodParameters(paymentMethodParameters).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to create payment method: %w", err)
	}

	if createPMResp == nil {
		return nil, fmt.Errorf("create payment method response is nil")
	}

	return createPMResp.ReferenceId, nil
}

func CardPaymentMethod(card *models.PaymentMethodCard) (*string, error) {
	paymentMethodParameters := *payment_method.
		NewPaymentMethodParameters(
			payment_method.PaymentMethodType("CARD"),
			payment_method.PaymentMethodReusability("ONE_TIME_USE"))

	successReturnUrl := "https://redirect.me/goodstuff"
	failureReturnUrl := "https://redirect.me/badstuff"

	cardParams := payment_method.CardParameters{
		Currency: "IDR",
		CardInformation: &payment_method.CardParametersCardInformation{
			CardNumber:     card.CardNumber,
			ExpiryMonth:    card.ExpiryMonth,
			ExpiryYear:     card.ExpiryYear,
			CardholderName: *payment_method.NewNullableString(card.CardholderName),
			Cvv:            *payment_method.NewNullableString(card.Cvv),
		},
		ChannelProperties: *payment_method.NewNullableCardChannelProperties(&payment_method.CardChannelProperties{
			SuccessReturnUrl: *payment_method.NewNullableString(&successReturnUrl),
			FailureReturnUrl: *payment_method.NewNullableString(&failureReturnUrl),
		}),
	}

	paymentMethodParameters.SetCard(cardParams)

	createPMResp, _, err := Client.PaymentMethodApi.
		CreatePaymentMethod(context.Background()).
		PaymentMethodParameters(paymentMethodParameters).
		Execute()

	return createPMResp.ReferenceId, err
}
