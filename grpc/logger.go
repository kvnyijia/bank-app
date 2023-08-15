package grpc

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg(">>> received a gRPC request")
	return result, err
}

type MyResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (writer *MyResponseWriter) WriteHeader(statusCode int) {
	writer.StatusCode = statusCode
	writer.ResponseWriter.WriteHeader(statusCode)
}

func (writer *MyResponseWriter) Write(body []byte) (int, error) {
	writer.Body = body
	return writer.ResponseWriter.Write(body)
}

func HttpHandler(handler http.Handler) http.Handler {
	// http.HandlerFunc is a type alias of the function signature
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		resWriter := &MyResponseWriter{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(resWriter, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if resWriter.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", resWriter.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", int(resWriter.StatusCode)).
			Str("status_text", http.StatusText(resWriter.StatusCode)).
			Dur("duration", duration).
			Msg(">>> received a http request")
	})
}
