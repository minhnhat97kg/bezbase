package docs

import "github.com/swaggo/swag"

// @title BezBase API
// @version 1.0
// @description This is the API documentation for BezBase application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.bezbase.com/support
// @contact.email support@bezbase.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func SwaggerInfo() {
	swag.Register(swag.Name, &swag.Spec{
		InfoInstanceName: "swagger",
		// SwaggerTemplate: docTemplate,
	})
}
