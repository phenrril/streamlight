package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"log"
	"opennode/internal/login/login_handler"
	"opennode/internal/pago/handler/input"
	"opennode/internal/pago/port"
	"opennode/internal/pago/usecase"
	"opennode/internal/renders/render_handler"
)

type Server struct {
	router       *gin.Engine
	oauth2Config *oauth2.Config
}

func NewServer() (*Server, error) {
	router := gin.New()
	return &Server{
		router: router,
	}, nil
}

func main() {
	//Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	//Load Server
	server, err := NewServer()
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	//Oauth2 server
	oauth2Config, err := login_handler.NewOauth2Config()
	if err != nil {
		log.Fatalf("Failed to get oauth2 config: %v", err)
	}

	oauthServer, err := login_handler.NewServer(oauth2Config)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	//Session Storage
	store := cookie.NewStore([]byte("SuperSecretVault"))
	server.router.Use(sessions.Sessions("auth-session", store))
	if err != nil {
		log.Fatalf("Failed to create session storage: %v", err)
	}

	portPago := port.NewPagoPort()
	useCase := usecase.NewPagoUsecase(portPago)
	pagoHandler := input.NewPagoHandler(useCase, store)

	//CSS and Js Handlers
	server.router.Static("/web/static/css", "./web/static/css")
	server.router.Static("/web/static/js", "./web/static/js")

	//Route Handlers

	server.router.GET("/login", oauthServer.LoginHandler)
	server.router.GET("/logout", oauthServer.LogoutHandler)
	server.router.GET("/callback", oauthServer.CallbackHandler)

	server.router.POST("/pago", pagoHandler.HandlePago)

	server.router.GET("/index" /*login_handler.IsAuthenticated(),*/, render_handler.Index)
	server.router.GET("/home", render_handler.HomeTemplate)
	server.router.GET("/about", render_handler.About)
	server.router.GET("/faq", render_handler.FaqTemplate)

	server.router.StaticFile("/web/static/img/logo2.png", "./web/static/img/logo2.png")
	server.router.StaticFile("/web/static/img/sinfondo.png", "./web/static/img/sinfondo.png")

	//Temporary Route Successful Payment
	server.router.GET("/exito", func(c *gin.Context) {

		session := sessions.Default(c)
		amount := session.Get("amount")
		currrency := session.Get("currency")
		description := session.Get("description")
		customerName := session.Get("customer_name")
		customerEmail := session.Get("customer_email")
		orderID := session.Get("order_id")
		callbackURL := session.Get("callback_url")
		successURL := session.Get("success_url")
		autoSettle := session.Get("auto_settle")
		splitToBtcBps := session.Get("split_to_btc_bps")
		hostedCheckoutURL := session.Get("hosted_checkout_url")

		c.Writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(c.Writer, "Amount: "+fmt.Sprint(amount))
		fmt.Fprintln(c.Writer, "Currency: "+fmt.Sprint(currrency))
		fmt.Fprintln(c.Writer, "Description: "+fmt.Sprint(description))
		fmt.Fprintln(c.Writer, "Customer Name: "+fmt.Sprint(customerName))
		fmt.Fprintln(c.Writer, "Customer Email: "+fmt.Sprint(customerEmail))
		fmt.Fprintln(c.Writer, "Order ID: "+fmt.Sprint(orderID))
		fmt.Fprintln(c.Writer, "Callback URL: "+fmt.Sprint(callbackURL))
		fmt.Fprintln(c.Writer, "Success URL: "+fmt.Sprint(successURL))
		fmt.Fprintln(c.Writer, "Auto Settle: "+fmt.Sprint(autoSettle))
		fmt.Fprintln(c.Writer, "Split To Btc Bps: "+fmt.Sprint(splitToBtcBps))
		fmt.Fprintln(c.Writer, "Hosted Checkout URL: "+fmt.Sprint(hostedCheckoutURL))
	})

	if err := server.router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
