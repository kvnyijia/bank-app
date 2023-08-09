package grpc

import (
	"context"

	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/util"
	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	// Hash the given password
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, ">>> failed to hash password: %s", err)
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.Internal, ">>> the username already exists: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, ">>> failed to create user: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return res, nil
}
