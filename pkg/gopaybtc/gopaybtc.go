package gopaybtc

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PaymentRequest struct {
	Amount float64 `json:"amount"`
}

type PaymentResponse struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

type AddressBalanceResponse struct {
	Address string  `json:"address"`
	Balance float64 `json:"balance"`
}

type Server struct {
	router             *gin.Engine
	blockonomicsClient *BlockonomicsClient
}

func NewServer(config Config) (*Server, error) {
	server := &Server{
		blockonomicsClient: NewBlockonomicsClient(config),
	}
	server.setUpRouter()

	return server, nil
}

// Start Run the server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) healthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func (server *Server) handlePayment(ctx *gin.Context) {
	amount, err := strconv.Atoi(ctx.Query("amount"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	callbackURL := ctx.Query("callback_url")
	if callbackURL == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(errors.New("callback_url parameter is missing")))
		return
	}

	resp, err := server.blockonomicsClient.CreatePaymentRequest(amount, callbackURL)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (server *Server) handleAddressInfo(ctx *gin.Context) {
	address := ctx.Param("address")
	if address == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(fmt.Errorf("address parameter is required")))
		return
	}
	addressInfo, err := server.blockonomicsClient.GetAddressInfo(address)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, addressInfo)
}

func (server *Server) handleTransactionInfo(ctx *gin.Context) {
	txid := ctx.Param("txid")
	if txid == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, errorResponse(fmt.Errorf("txid parameter is required")))
		return
	}
	txInfo, err := server.blockonomicsClient.GetTransactionInfo(txid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, txInfo)
}

func (server *Server) setUpRouter() {
	router := gin.Default()

	router.POST("/v1/payment", server.handlePayment)
	router.GET("/v1/health", server.healthCheck)
	router.GET("/v1/transaction/:txid", server.handleTransactionInfo)
	router.GET("/v1/address/:address", server.handleAddressInfo)

	server.router = router

}

// errors
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
