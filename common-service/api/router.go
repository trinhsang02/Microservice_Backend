package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the API routes
func SetupRouter(h *Handler) *gin.Engine {
	r := gin.Default()

	// Add CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Or specify your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// KYC endpoints
	r.POST("/kyc", h.SubmitKYC)
	r.GET("/kyc/citizen/:citizenID", h.GetKYCByCitizenID)
	// r.GET("/kyc/wallet/:walletAddress", h.GetKYCByWalletAddress)
	r.PUT("/kyc", h.UpdateKYC)

	// Event endpoints
	r.GET("/events/:netId/:contractAddress/:eventType", h.GetEvents)
	r.GET("/event/:eventType/:hex", h.GetEventByInfo)
	r.GET("/leaves/:netId/:contractAddress", h.GetLeaves)

	return r
}
