package grpc

import (
	"context"

	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/util"
	"github.com/lib/pq"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(req)
	if violations != nil {
		return nil, invalidArgErr(violations)
	}
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

func validateCreateUserRequest(req *pb.CreateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var arr []*errdetails.BadRequest_FieldViolation
	if err := ValidateUsername(req.GetUsername()); err != nil {
		arr = append(arr, fieldViolation("username", err))
	}
	if err := ValidatePassword(req.GetPassword()); err != nil {
		arr = append(arr, fieldViolation("password", err))
	}
	if err := ValidateFullname(req.GetFullName()); err != nil {
		arr = append(arr, fieldViolation("full_name", err))
	}
	if err := ValidateEmail(req.GetEmail()); err != nil {
		arr = append(arr, fieldViolation("email", err))
	}
	return arr
}
