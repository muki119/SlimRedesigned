package routes

import (
	"net/http"
)

func (ur *UserRoutes) GetUserRoutes() *http.ServeMux {
	userMux := http.NewServeMux()
	userMux.HandleFunc("GET /users/me", ur.MiddleWare.ErrorHandler(ur.MiddleWare.UserAuthMiddleware(ur.UserControllers.GetUserProfile))) // be wrapped in error handler and auth middleware
	userMux.HandleFunc("PUT /users/me", ur.MiddleWare.ErrorHandler(ur.MiddleWare.UserAuthMiddleware(ur.UserControllers.UpdateProfile)))  // be wrapped in error handler and auth middleware
	userMux.HandleFunc("GET /users/:username", ur.MiddleWare.ErrorHandler(ur.UserControllers.GetProfile))
	return userMux
}
