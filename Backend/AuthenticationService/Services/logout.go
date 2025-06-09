package Services

import (
	"github.com/golang-jwt/jwt/v5"
)

func LogoutService(incommingToken *jwt.Token) error { // set last active time to current time and invalidate user session

	// Clear authentication tokens or cookies
	// Return success or error response
	return nil
}
