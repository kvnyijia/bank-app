package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/kvnyijia/bank-app/token"
	"google.golang.org/grpc/metadata"
)

const (
	authField           = "authorization"
	authorizationBearer = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	mdata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf(">>> missing metadata")
	}
	values := mdata.Get(authField)
	if len(values) == 0 {
		return nil, fmt.Errorf(">>> missing authorization header")
	}

	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf(">>> invalid authorization header format")
	}

	authType := strings.ToLower(fields[0])
	if authType != authorizationBearer {
		return nil, fmt.Errorf(">>> unsupported authorization type: %s", authType)
	}

	accessToken := fields[1]
	payload, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf(">>> invalid access token: %s", err)
	}
	return payload, nil
}