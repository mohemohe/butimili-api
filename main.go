package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mohemohe/butimili-api/configs"
	"github.com/mohemohe/butimili-api/controllers"
	"github.com/mohemohe/butimili-api/controllers/api/v1"
	_ "github.com/mohemohe/butimili-api/docs"
	"github.com/mohemohe/butimili-api/middlewares"
	"github.com/mohemohe/butimili-api/models"
	"github.com/mohemohe/echoHelper"
	"github.com/swaggo/echo-swagger"
)

// @Title butimili API
// @Version v1
// @SecurityDefinitions.apikey AccessToken
// @In header
// @Name Authorization
// @AuthorizationURL /api/v1/auth
func main() {
	eh := echoHelper.New(echo.New(), echoHelper.WithCustomMiddleware([]echo.MiddlewareFunc{
		middleware.Recover(),
		middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `time="${time_rfc3339}" level=info msg="echo" remote_ip="${remote_ip}" ` +
				`host="${host}" method="${method}" uri="${uri}" status="${status} ` +
				`latency="${latency}" latency_human="${latency_human}" ` +
				`bytes_in="${bytes_in}" bytes_out="${bytes_out}"` + "\n",
		}),
		//middleware.StaticWithConfig(middleware.StaticConfig{
		//	Root: "frontend/dist",
		//	Skipper: func(c echo.Context) bool {
		//		return strings.HasPrefix(c.Request().URL.Path, "/api") || strings.HasPrefix(c.Request().URL.Path, "/public")
		//	},
		//}),
		middleware.CORS(),
	}))
	if configs.GetEnv().Echo.Env == "debug" {
		eh.Echo().Logger.SetLevel(0)
	}

	eh.RegisterRoutes([]echoHelper.Route{
		{echo.GET, "/swagger", controllers.RedirectSwagger, nil},
		{echo.GET, "/swagger/", controllers.RedirectSwagger, nil},
		{echo.GET, "/swagger/*", echoSwagger.WrapHandler, nil},

		{echo.GET, "/api/v1/version", v1.GetVersion, nil},
		{echo.GET, "/api/v1/auth", v1.GetAuth, &[]echo.MiddlewareFunc{
			middlewares.AuthMiddleware,
		}},
		{echo.POST, "/api/v1/auth", v1.PostAuth, nil},
		{echo.POST, "/api/v1/user", v1.PostUser, nil},
		{echo.DELETE, "/api/v1/user", v1.DeleteUser, &[]echo.MiddlewareFunc{
			middlewares.AuthMiddleware,
		}},
		{echo.GET, "/api/v1/butimili", v1.ListButimiliText, &[]echo.MiddlewareFunc{
			middlewares.AuthMiddleware,
		}},
		{echo.GET, "/api/v1/butimili/list", v1.ListButimili, &[]echo.MiddlewareFunc{
			middlewares.AuthMiddleware,
		}},
		{echo.PUT, "/api/v1/butimili/list", v1.PutButimili, &[]echo.MiddlewareFunc{
			middlewares.AuthMiddleware,
		}},
		{echo.DELETE, "/api/v1/butimili/list/:screenName", v1.DeleteButimili, &[]echo.MiddlewareFunc{
			middlewares.AuthMiddleware,
		}},
		{echo.GET, "/api/v1/butimili/raw", v1.GetButimiliText, nil},
	})

	go models.InitDB()

	//eh.Serve()
	eh.Echo().Logger.Fatal(eh.Echo().Start(configs.GetEnv().Echo.Address))
}
