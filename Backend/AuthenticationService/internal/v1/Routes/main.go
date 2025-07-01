package Routes

import (
	"v1/Controllers"
	"v1/Middleware"
)

type Routes struct {
	RouteControllers *Controllers.Controllers
	Middleware       *Middleware.Middleware
}
