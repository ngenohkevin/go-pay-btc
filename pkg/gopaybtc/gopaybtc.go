package gopaybtc

import (
	"github.com/gin-gonic/gin"
	"log"
)

// PaymentGateway represents the payment gateway object
type PaymentGateway struct {
	Router *gin.Engine
}

// NewPaymentGateway creates a new payment gateway object
//func NewPaymentGateway() *PaymentGateway {
//	//pg := &PaymentGateway{}
//	//pg.Router = gin.Default()
//	//pg.setRoutes()
//	//return pg
//}

// setRoutes sets up the routes for the payment gateway
func (pg *PaymentGateway) setRoutes() {
	//pg.Router.POST("/payment", pg.handlePayment)
	//pg.Router.GET("/status/:id", pg.handleStatus)
}

// handlePayment handles the payment request
func (pg *PaymentGateway) handlePayment(c *gin.Context) {
	// Handle payment request here
}

// handleStatus handles the payment status request
func (pg *PaymentGateway) handleStatus(c *gin.Context) {
	// Handle payment status request here
}

// Start starts the payment gateway
func (pg *PaymentGateway) Start(port string) {
	err := pg.Router.Run(port)
	if err != nil {
		log.Fatal("Error starting payment gateway: ", err)
	}
}
