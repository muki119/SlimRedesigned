package middleware

import (
	"v1/Helpers/token"
)

type Middleware struct {
	TokenHelper *token.Token
}
