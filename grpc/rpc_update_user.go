package grpc

import (
	"context"
	"database/sql"
	"time"

	db "github.com/kvnyijia/bank-app/db/sqlc"
	"github.com/kvnyijia/bank-app/pb"
	"github.com/kvnyijia/bank-app/util"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO: add auth

	violations := validateUpdateUserRequest(req)
	if violations != nil {
		return nil, invalidArgErr(violations)
	}

	arg := db.UpdateUserParams{
		Username: req.GetUsername(),
		FullName: sql.NullString{
			String: req.GetFullName(),
			Valid:  req.FullName != nil,
		},
		Email: sql.NullString{
			String: req.GetEmail(),
			Valid:  req.Email != nil,
		},
	}

	if req.Password != nil {
		// Hash the given password
		hashedPassword, err := util.HashPassword(req.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, ">>> failed to hash password: %s", err)
		}
		arg.HashedPassword = sql.NullString{
			String: hashedPassword,
			Valid:  true,
		}
		arg.PasswordChangedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUser(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, ">>> the username does not exists: %s", err)
		}
		return nil, status.Errorf(codes.Internal, ">>> failed to Update user: %s", err)
	}

	res := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return res, nil
}

func validateUpdateUserRequest(req *pb.UpdateUserRequest) []*errdetails.BadRequest_FieldViolation {
	var arr []*errdetails.BadRequest_FieldViolation
	if err := ValidateUsername(req.GetUsername()); err != nil {
		arr = append(arr, fieldViolation("username", err))
	}
	if req.Password != nil {
		if err := ValidatePassword(req.GetPassword()); err != nil {
			arr = append(arr, fieldViolation("password", err))
		}
	}
	if req.FullName != nil {
		if err := ValidateFullname(req.GetFullName()); err != nil {
			arr = append(arr, fieldViolation("full_name", err))
		}
	}
	if req.Email != nil {
		if err := ValidateEmail(req.GetEmail()); err != nil {
			arr = append(arr, fieldViolation("email", err))
		}
	}
	return arr
}
