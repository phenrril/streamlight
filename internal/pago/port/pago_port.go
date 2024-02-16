package port

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"opennode/internal/pago/model"
	"strings"
)

type PagoPort interface {
	RealizarPago(pago model.Pago) (model.Pago, error)
}

type ResponseData struct {
	HostedCheckoutURL string `json:"hosted_checkout_url"`
}

type Response struct {
	Data ResponseData `json:"data"`
}

type PagoPortImpl struct{}

func (p PagoPortImpl) RealizarPago(pago model.Pago) (model.Pago, error) {
	// Crear la solicitud a la API de pagos
	url := "https://api.opennode.com/v1/charges"
	payload := strings.NewReader(fmt.Sprintf(
		`{"amount":%d,
				"currency":"%s",
				"description":"%s",
				"customer_name":"%s",
				"customer_email":"%s",
				"order_id":"%s",
				"callback_url":"%s",
				"success_url":"%s",
				"auto_settle":%t,
				"split_to_btc_bps":%d}`,
		pago.Amount,
		pago.Currency,
		pago.Description,
		pago.CustomerName,
		pago.CustomerEmail,
		pago.OrderID,
		pago.CallbackURL,
		pago.SuccessURL,
		pago.AutoSettle,
		pago.SplitToBtcBps,
	))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "904ade28-657e-449e-837a-1a34c08ce708")

	// Enviar la solicitud y obtener la respuesta
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return pago, fmt.Errorf("Error al realizar la solicitud a la API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error al cerrar el cuerpo de la respuesta de la API:", err)
		}
	}(res.Body)

	// Leer y devolver el cuerpo de la respuesta
	body, _ := io.ReadAll(res.Body)

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Error al deserializar la respuesta de la API:", err)
		return pago, err
	}

	if response.Data.HostedCheckoutURL == "" {
		return pago, fmt.Errorf("La respuesta de la API no contiene un valor para HostedCheckoutURL")
	}

	pago.HostedCheckoutURL = response.Data.HostedCheckoutURL
	return pago, nil
}

func NewPagoPort() PagoPort {
	return PagoPortImpl{}
}
