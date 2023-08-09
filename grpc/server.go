package grpc

import (
	"fmt"

	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/token"
	"github.com/kvnyijia/bank-app/util"
)

// Server serves gRPC reqs for our banking service
type Server struct {
	pb.UnimplementedBankAppServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// Creates a new gRPC server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) // Could also be NewJWTMaker()
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
