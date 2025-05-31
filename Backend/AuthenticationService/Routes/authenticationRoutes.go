package Routes

import (
	"net/http"
	"v1/Controllers"
	"v1/Middleware"
)

var AuthRouter = http.NewServeMux()

func Initial() {
	AuthRouter.HandleFunc("POST /register", Middleware.ErrorHandler(Controllers.RegisterHandler))
	AuthRouter.HandleFunc("POST /login", Middleware.ErrorHandler(Controllers.LoginHandler))
	AuthRouter.HandleFunc("DELETE /logout", Middleware.ErrorHandler(Middleware.CheckUserLoggedIn(Controllers.LogoutHandler)))
}
