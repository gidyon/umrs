package contract

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type contractAuth struct {
	contractID string
}

// NewledgerContractAuth creates a PerRPCCredentials auth contract
func NewledgerContractAuth(contractID string) credentials.PerRPCCredentials {
	return &contractAuth{contractID}
}

func (c *contractAuth) GetRequestMetadata(
	ctx context.Context, strs ...string,
) (map[string]string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.NotFound, "no metadata found in context")
	}
	reqMD := make(map[string]string, 0)
	for key, value := range md {
		if len(value) != 0 {
			reqMD[key] = value[0]
		}
	}
	reqMD["contract_id"] = c.contractID
	return reqMD, nil
}

func (*contractAuth) RequireTransportSecurity() bool {
	return true
}
