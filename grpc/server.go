package grpc

import (
	"fmt"

	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/token"
	"github.com/kvnyijia/bank-app/util"
	"github.com/kvnyijia/bank-app/worker"
)

// Server serves gRPC reqs for our banking service
type Server struct {
	pb.UnimplementedBankAppServer
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// Creates a new gRPC server
func NewServer(config util.Config, store db.Store, taskDist worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) // Could also be NewJWTMaker()
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDist,
	}

	return server, nil
}
