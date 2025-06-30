package Routes

import (
	"net/http"
	"v1/Controllers"
	"v1/Middleware"
)

func InitialiseRoutes(RouteControllers *Controllers.Controllers) *http.ServeMux {
	var AuthRouter = http.NewServeMux()
	AuthRouter.HandleFunc("POST /register", Middleware.ErrorHandler(RouteControllers.RegisterHandler))
	AuthRouter.HandleFunc("POST /login", Middleware.ErrorHandler(RouteControllers.LoginHandler))
	AuthRouter.HandleFunc("GET /token", Middleware.ErrorHandler(Middleware.CheckUserLoggedIn(RouteControllers.TokenHandler)))
	AuthRouter.HandleFunc("DELETE /logout", Middleware.ErrorHandler(Middleware.CheckUserLoggedIn(RouteControllers.LogoutHandler)))
	return AuthRouter
}
