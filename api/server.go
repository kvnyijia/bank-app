package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/kvnyijia/bank-app/db/sqlc"
)

// Server serves HTTP reqs for our banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// Add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start to run the HTTP server on an address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
