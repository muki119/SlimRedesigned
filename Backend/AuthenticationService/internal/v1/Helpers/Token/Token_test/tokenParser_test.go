package Token_test

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"testing"
	"time"
	"v1/Config"
	"v1/Helpers/Token"
	"v1/Utils"
)

func TestParseRefreshToken(t *testing.T) {
	userId := uuid.New().String()
	issuer := "/test"

	TokenBlocklistDB := &Config.RedisConfig{
		Addr:     Utils.MustGetEnv("REDIS_HOST"),
		Password: Utils.MustGetEnv("REDIS_PASSWORD"),
		DB:       1,
	}
	TokenConfig := &Token.HelperTokenConfig{
		SecretKey:  "JWT_SECRET_KEY",
		PrivateKey: "JWT_PRIVATE_KEY",
	}
	TokenHelper := TokenConfig.CreateTokenService(TokenBlocklistDB.ConnectToDatabase())
	t.Run("AccessTokenShouldFailAsRefreshToken", func(t *testing.T) {
		accessToken, err := TokenHelper.CreateAccessToken(userId, issuer)
		if err != nil {
			t.Fatalf("Failed to create access token: %v", err)
		}

		_, err = TokenHelper.ParseRefreshToken(accessToken)
		if err == nil {
			t.Error("Expected error when parsing access token as refresh token, but got none")
		}
	})

	t.Run("ValidRefreshTokenShouldParse", func(t *testing.T) {
		refreshToken, err := TokenHelper.CreateLoginRefreshToken(userId)
		if err != nil {
			t.Fatalf("Failed to create refresh token: %v", err)
		}

		parsedToken, err := TokenHelper.ParseRefreshToken(refreshToken)
		if err != nil {
			t.Fatalf("Failed to parse valid refresh token: %v", err)
		}

		// Verify the parsed token contains expected claims
		claims := parsedToken.Claims.(jwt.MapClaims)
		if claims["sub"] != userId {
			t.Errorf("Expected user ID %s, got %v", userId, claims["sub"])
		}

		// Verify token type
		if tokenType, ok := claims["typ"]; ok {
			if tokenType != "refresh" {
				t.Errorf("Expected token type 'refresh', got %v", tokenType)
			}
		}
	})

	t.Run("RefreshTokenReturnErrorIfNoRefreshToken", func(t *testing.T) {
		_, err := TokenHelper.ParseRefreshToken("")
		if err == nil {
			t.Error("Expected error when parsing refresh token, but got none")
		}
	})

	t.Run("ExpiredRefreshTokenShouldReturn", func(t *testing.T) {
		expiredClaims := jwt.RegisteredClaims{
			Subject:   userId,
			Audience:  jwt.ClaimStrings{":("},
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(-1, 0, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
			ID:        userId,
		}
		expiredToken, err := TokenHelper.CreateRefreshTokenFromClaims(expiredClaims)
		if err != nil {
			t.Fatalf("Failed to create refresh token: %v", err)
		}
		_, err = TokenHelper.ParseRefreshToken(expiredToken)
		if err == nil {
			t.Error("Expected error when parsing refresh token, but got none")
		}
		if !errors.Is(err, jwt.ErrTokenInvalidClaims) {
			t.Error("Expecter expired token error ", "got ", err)
		}
	})

}
