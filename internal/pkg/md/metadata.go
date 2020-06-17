package md

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// AddFromCtx extracts metadata in ctx and create an outgoing context with the MD attached
func AddFromCtx(ctx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	return metadata.NewOutgoingContext(ctx, md)
}
