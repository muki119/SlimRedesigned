package Routes

import (
	"AuthenticationService/v1/Controllers"
	"AuthenticationService/v1/Middleware"
)

type Routes struct {
	RouteControllers *Controllers.Controllers
	Middleware       *Middleware.Middleware
}
