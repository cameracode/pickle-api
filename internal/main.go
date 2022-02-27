package main

import (
	"fmt"
	"log"
	"net/http"

	"pickle-api/pkg/swagger/server/restapi"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"

	"pickle-api/pkg/swagger/server/restapi/operations"
)

func main() {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPickleAPIAPI(swaggerSpec)
	server := restapi.NewServer(api)

	defer func() {
		if err := server.Shutdown(); err != nil {
			// error handler
			log.Fatalln(err)
		}
	}()

	server.Port = 8080

	api.CheckHealthHandler = operations.CheckHealthHandlerFunc(Health)

	api.GetHelloUserHandler = operations.GetHelloUserHandlerFunc(GetHelloUser)

	api.GetPickleNameHandler = operations.GetPickleNameHandlerFunc(GetPickleByName)

	// Start listening server
	if err := server.Serve(); err != nil {
		log.Fatalln(err)
	}
}

// Health route returns 200 - OK
func Health(user operations.CheckHealthParams) middleware.Responder {
	return operations.NewCheckHealthOK().WithPayload("OK")
}

// GetHelloUser returns Hello + your name
func GetHelloUser(user operations.GetHelloUserParams) middleware.Responder {
	return operations.NewGetHelloUserOK().WithPayload("Hello " + user.User + "!")
}

// GetPickleByName returns a Pickle Rick in PNG format
func GetPickleByName(pickle operations.GetPickleNameParams) middleware.Responder {
	var URL string
	if pickle.Name != "" {
		//URL = "https://github.com/cameracode/pickle-api/assets/picklericks/raw/" + pickle.Name + ".png"
		URL = "https://raw.githubusercontent.com/cameracode/ricksofpickle/Develop/" + pickle.Name + ".png"
	} else {
		// by default we return Oscar Arakaki's cell-shaded 3D Pickle Rick
		// Art Credit: https://oki93.artstation.com/projects/K5x2y
		//URL = "https://github.com/cameracode/pickle-api/assets/picklericks/raw/arakaki-picklerick.png"
		// hardcode URL cuz im a noob
		URL = "https://raw.githubusercontent.com/cameracode/ricksofpickle/Develop/arakaki-picklerick.png"
	}

	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error")
	}

	return operations.NewGetPickleNameOK().WithPayload(response.Body)
}
