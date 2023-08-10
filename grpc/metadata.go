package grpc

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	UserAgentField  = "grpcgateway-user-agent"
	UserAgentField2 = "user-agent"
	ClientIpField   = "x-forwarded-host"
)

type MetaData struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(ctx context.Context) *MetaData {
	mdata := &MetaData{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(UserAgentField); len(userAgents) > 0 {
			mdata.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(UserAgentField2); len(userAgents) > 0 {
			mdata.UserAgent = userAgents[0]
		}

		if ClientIps := md.Get(ClientIpField); len(ClientIps) > 0 {
			mdata.ClientIp = ClientIps[0]
		}
	}

	if peerData, ok := peer.FromContext(ctx); ok {
		mdata.ClientIp = peerData.Addr.String()
	}
	return mdata
}
