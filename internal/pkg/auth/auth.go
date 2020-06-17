package auth

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/umrs/internal/pkg/errs"
	"github.com/gidyon/umrs/pkg/api/ledger"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
)

// PatientGroup is group for a patient user
func PatientGroup() string {
	return ledger.Actor_PATIENT.String()
}

// AdminGroup is group for admin user
func AdminGroup() string {
	return ledger.Actor_ADMIN.String()
}

// SuperAdminGroup is group for super admin user
func SuperAdminGroup() string {
	return ledger.Actor_ADMIN.String()
}

// InsuranceGroup is group for insurance user
func InsuranceGroup() string {
	return ledger.Actor_INSURANCE.String()
}

// HospitalGroup is hospital user
func HospitalGroup() string {
	return ledger.Actor_HOSPITAL.String()
}

// Interface is a generic authentication and authorization API
type Interface interface {
	AuthenticateRequest(context.Context) error
	AuthorizeActor(ctx context.Context, actorID string) (*Payload, error)
	AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*Payload, error)
	AuthorizeStrict(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error)
	GenToken(context.Context, *Payload, int64) (string, error)
}

type authAPI struct {
	signingKey string
}

// NewAPI creates new auth API with given signing key
func NewAPI(signingKey string) (Interface, error) {
	api := &authAPI{signingKey}
	return api, nil
}

func (api *authAPI) AuthenticateRequest(ctx context.Context) error {
	_, err := api.ParseFromCtx(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (api *authAPI) AuthorizeActor(ctx context.Context, actorID string) (*Payload, error) {
	claims, err := api.ParseFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	if claims.ID != actorID {
		return nil, errs.TokenCredentialNotMatching("ID")
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*Payload, error) {
	claims, err := api.ParseFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	err = matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) AuthorizeStrict(ctx context.Context, actorID string, allowedGroups ...string) (*Payload, error) {
	claims, err := api.ParseFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	err = matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	if claims.ID != actorID {
		return nil, err
	}

	return claims.Payload, nil
}

func (api *authAPI) GenToken(ctx context.Context, payload *Payload, expires int64) (string, error) {
	return api.genToken(ctx, payload, expires)
}

// AddMD adds metadata to token
func (api *authAPI) AddMD(ctx context.Context, actorID, group string) context.Context {
	payload := &Payload{
		ID:           actorID,
		FirstName:    randomdata.FirstName(randomdata.Male),
		LastName:     randomdata.LastName(),
		EmailAddress: randomdata.Email(),
		Group:        group,
		Label:        "",
	}
	token, err := api.genToken(ctx, payload, 0)
	if err != nil {
		panic(err)
	}

	return addTokenMD(ctx, token)
}

// ParseToken parses a jwt token and return claims or error if token is invalid
func (api *authAPI) ParseToken(tokenString string) (claims *Claims, err error) {
	// Handling any panic is good trust me!
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return api.signingKey, nil
		},
	)
	if err != nil {
		return nil, status.Errorf(
			codes.Unauthenticated, "failed to parse token with claims: %v", err,
		)
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "the token is not valid")
	}
	return claims, nil
}

// ParseFromCtx jwt token from context
func (api *authAPI) ParseFromCtx(ctx context.Context) (*Claims, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(
			codes.Internal, "failed to get Bearer from authorization header: %v", err,
		)
	}

	return api.ParseToken(token)
}

func addTokenMD(ctx context.Context, token string) context.Context {
	return metadata.NewIncomingContext(
		ctx, metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token)),
	)
}

func matchGroup(claimGroup string, allowedGroups []string) error {
	for _, group := range allowedGroups {
		if claimGroup == group {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied, "permission denied for group %s", claimGroup)
}

func (api *authAPI) genToken(
	ctx context.Context, payload *Payload, expires int64,
) (tokenStr string, err error) {
	// Handling any panic is good trust me!
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("%v", err2)
		}
	}()

	token := jwt.NewWithClaims(signingMethod, Claims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
			Issuer:    "umrs",
		},
	})

	// Generate the token
	return token.SignedString(api.signingKey)
}
