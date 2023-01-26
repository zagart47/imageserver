package internal

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync/atomic"
)

type ServerMiddleware struct {
	loadConn    int64
	getListConn int64
}

func (sm *ServerMiddleware) Interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	switch info.FullMethod {
	case "/proto.FileService/Upload",
		"/proto.FileService/Download":
		if atomic.LoadInt64(&sm.loadConn) > 10 {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected, please retry later", info.FullMethod)
		}
		atomic.AddInt64(&sm.loadConn, 1)
		defer atomic.AddInt64(&sm.loadConn, -1)
		return handler(ctx, req)
	case "/proto.FileService/GetFiles":
		if atomic.LoadInt64(&sm.getListConn) > 100 {
			return nil, status.Errorf(codes.ResourceExhausted, "%s is rejected, please retry later", info.FullMethod)
		}
		atomic.AddInt64(&sm.getListConn, 1)
		defer atomic.AddInt64(&sm.getListConn, -1)
		return handler(ctx, req)
	default:
		return handler(ctx, req)
	}
}
