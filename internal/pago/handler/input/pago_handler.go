package input

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
	"opennode/internal/pago/model"
	"opennode/internal/pago/usecase"
)

type PagoHandler struct {
	pagoUsecase usecase.PagoUsecase
	store       cookie.Store
}

func (h PagoHandler) savePagoInSession(c *gin.Context, pago model.Pago) {
	session := sessions.Default(c)

	session.Set("amount", pago.Amount)
	session.Set("currency", pago.Currency)
	session.Set("description", pago.Description)
	session.Set("customer_name", pago.CustomerName)
	session.Set("customer_email", pago.CustomerEmail)
	session.Set("order_id", pago.OrderID)
	session.Set("callback_url", pago.CallbackURL)
	session.Set("success_url", pago.SuccessURL)
	session.Set("auto_settle", pago.AutoSettle)
	session.Set("split_to_btc_bps", pago.SplitToBtcBps)
	session.Set("hosted_checkout_url", pago.HostedCheckoutURL)
	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func NewPagoHandler(u usecase.PagoUsecase, store cookie.Store) PagoHandler {
	return PagoHandler{
		pagoUsecase: u,
		store:       store,
	}
}

func (h PagoHandler) HandlePago(c *gin.Context) {
	if c.Request.Method == "GET" {
		// form HTML
	} else if c.Request.Method == "POST" {
		var body model.Pago
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		pago := model.NewPago(
			body.Amount,
			body.CustomerName,
			body.CustomerEmail,
			body.HostedCheckoutURL,
		)
		pago, err = h.pagoUsecase.RealizarPago(pago)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		h.savePagoInSession(c, pago)
		c.Redirect(http.StatusSeeOther, "/exito")
	}
}
