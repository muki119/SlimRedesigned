package Routes

import (
	"net/http"
)

func (routes *Routes) GetAuthRoutes() *http.ServeMux {
	var AuthRouter = http.NewServeMux()
	AuthRouter.HandleFunc("POST /register", routes.Middleware.ErrorHandler(routes.RouteControllers.RegisterHandler))
	AuthRouter.HandleFunc("POST /login", routes.Middleware.ErrorHandler(routes.RouteControllers.LoginHandler))
	AuthRouter.HandleFunc("GET /token", routes.Middleware.ErrorHandler(routes.Middleware.CheckUserLoggedIn(routes.RouteControllers.TokenHandler)))
	AuthRouter.HandleFunc("DELETE /logout", routes.Middleware.ErrorHandler(routes.Middleware.CheckUserLoggedIn(routes.RouteControllers.LogoutHandler)))
	return AuthRouter
}
