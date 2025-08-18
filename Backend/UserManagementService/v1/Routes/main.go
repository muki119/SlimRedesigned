package routes

import (
	controllers "v1/Controllers"
	middleware "v1/Middleware"
)

type UserRoutes struct {
	UserControllers *controllers.UserControllers
	MiddleWare      *middleware.Middleware
}
