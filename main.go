// main.go

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	// Set the router as the default one provided by Gin
	router = gin.Default()

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	initializeRoutes()

	// Start serving the application
	router.Run()
}

// Render one of HTML, JSON or CSV based on the 'Accept' header of the request
// If the header doesn't specify this, HTML is rendered, provided that
// the template name is present
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	case "application/xml":
		// Respond with XML
		c.XML(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

func initializeRoutes() {

	// Handle the index route
	router.GET("/", showIndexPage)

	router.GET("/applications", getApplicationsPage)

	// Group application related routes together
	applicationRoutes := router.Group("/application")
	{
		// Handle GET requests at /article/view/some_article_id
		// applicationRoutes.GET("/view/:application_id", getApplication)

		// Handle the GET requests at /application/create
		// Show the application creation page
		// Ensure that the user is logged in by using the middleware
		applicationRoutes.GET("/create", showApplicationCreationPage)

		// Handle POST requests at /article/create
		// Ensure that the user is logged in by using the middleware
		applicationRoutes.POST("/create", createApplication)
	}
}
