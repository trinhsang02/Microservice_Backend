package api

import (

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	// KYC endpoints
	r.POST("/kyc", h.SubmitKYC)
	r.GET("/kyc/citizen/:citizenID", h.GetKYCByCitizenID)
	// r.GET("/kyc/wallet/:walletAddress", h.GetKYCByWalletAddress)
	r.PUT("/kyc", h.UpdateKYC)

	return r
}
