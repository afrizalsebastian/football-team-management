package helper

import (
	"context"

	"github.com/afrizalsebastian/football-team-management/application/dto"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secret string, claims *dto.AdminClaims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func VerifyToken(tokenString string, secret string) (*dto.AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.AdminClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*dto.AdminClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

const ClaimsContextKey = "admin-info"

func SetContextValueClaims(ctx context.Context, value *dto.AdminClaims) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, value)
}

func GetContextValueClaims(ctx context.Context) *dto.AdminClaims {
	claims, ok := ctx.Value(ClaimsContextKey).(*dto.AdminClaims)
	if !ok {
		return nil
	}
	return claims
}
