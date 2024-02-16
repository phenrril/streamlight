package model

const (
	DefaultCurrency      = "BTC"
	DefaultAutoSettle    = false
	DefaultSplitToBtcBps = 10
	DefaultDescription   = "Opennode Payment"
	CallbackURL          = "http://localhost:8080/callback"
	SuccessURL           = "http://localhost:8080/success"
	FakeOrderId          = "fake_order_id"
)

type Pago struct {
	Amount            int    `json:"amount"`
	Currency          string `json:"currency"`
	Description       string `json:"description"`
	CustomerName      string `json:"customer_name"`
	CustomerEmail     string `json:"customer_email"`
	OrderID           string `json:"order_id"`
	CallbackURL       string `json:"callback_url"`
	SuccessURL        string `json:"success_url"`
	AutoSettle        bool   `json:"auto_settle"`
	SplitToBtcBps     int    `json:"split_to_btc_bps"`
	HostedCheckoutURL string `json:"hosted_checkout_url"`
}

func NewPago(
	amount int,
	customerName,
	customerEmail,
	hostedCheckoutURL string,
) Pago {
	return Pago{
		Amount:            amount,
		Currency:          DefaultCurrency,
		Description:       DefaultDescription,
		CustomerName:      customerName,
		CustomerEmail:     customerEmail,
		OrderID:           FakeOrderId,
		CallbackURL:       CallbackURL,
		SuccessURL:        SuccessURL,
		AutoSettle:        DefaultAutoSettle,
		SplitToBtcBps:     DefaultSplitToBtcBps,
		HostedCheckoutURL: hostedCheckoutURL,
	}
}
