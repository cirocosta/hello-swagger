// main declares the CLI that spins up the server of
// our API.
// It takes some arguments, validates if they're valid
// and match the expected type and then intiialize the
// server.
package main

import (
	"log"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/cirocosta/hello-swagger/swagger/models"
	"github.com/cirocosta/hello-swagger/swagger/restapi"
	"github.com/cirocosta/hello-swagger/swagger/restapi/operations"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
)

// cliArgs defines the configuration that the CLI
// expects. By using a struct we can very easily
// aggregate them into an object and check what are
// the expected types.
// If we need to mock this later it's just a matter
// of reusing the struct.
type cliArgs struct {
	Port int `arg:"-p,help:port to listen to"`
}

var (
	// args is a reference to an instantiation of
	// the configuration that the CLI expects but
	// with some values set.
	// By setting some values in advance we provide
	// default values that the user might provide
	// or not.
	args = &cliArgs{
		Port: 8080,
	}
)

// getHostnameHandler implements the handler that
// takes a set of parameters as described in swagger.yml
// and then produces a response.
// This response might be an error of a succesfull response.
//	-	In case of failure we create the payload
//		that would indicate the failure.
//	-	In case of success, the payload the indicates
//		the success with the hostname.
func getHostnameHandler(params operations.GetHostnameParams) middleware.Responder {
	payload, err := os.Hostname()

	if err != nil {
		errPayload := &models.Error{
			Code:    500,
			Message: swag.String("failed to retrieve hostname"),
		}

		return operations.
			NewGetHostnameDefault(500).
			WithPayload(errPayload)
	}

	return operations.NewGetHostnameOK().WithPayload(payload)
}

// main performs the main routine of the application:
//	1.	parses the args;
//	2.	analyzes the declaration of the API
//	3.	sets the implementation of the handlers
//	4.	listens on the port we want
func main() {
	arg.MustParse(args)

	// Load the JSON that corresponds to our swagger.yml
	// api definition.
	// This JSON is hardcoded as part of the generated code
	// that go-swagger creates.
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	// Load a dummy object that servers as an interface
	// that allows us to implement the API specification.
	api := operations.NewHelloAPI(swaggerSpec)

	// Create the REST api server that will make use of
	// the object that will container our handler implementations.
	server := restapi.NewServer(api)
	defer server.Shutdown()

	// Configure the server port
	server.Port = args.Port

	// Add our handler implementation
	api.GetHostnameHandler = operations.GetHostnameHandlerFunc(
		getHostnameHandler)

	// Let it run
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}
