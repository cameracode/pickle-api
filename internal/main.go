package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"net/http"
	"time"

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

	api.GetPicklesHandler = operations.GetGophersHandlerFunc(GetGophers)

	api.GetPickleRandomHandler = operations.GetPickleRandomHandlerFunc(GetPickleRandom)

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
		log.Fatalf("failed to open image: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		srcImage, _ := getPickleGopherError("Oops, Error")
		return operations.NewGetPickleNameOK().WithPayload(convertImgToIoCloser(srcImage))
	}

	srcImage, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatalf("failed to Decode image: %v", err)
	}

	if pickle.Size != nil {
		srcImage = resizeImage(srcImage, *pickle.Size)
	}

	return operations.NewGetPickleNameOK().WithPayload(convertImgToIoCloser(srcImage))
}

/*
Display Pickle List with optional filter
*/
func GetPickles(pickle operations.GetPicklesParams) middleware.Responder {
	picklesList := GetPicklesList()

	if pickle.Name != nil {
		var arr []*models.Pickle
		for key, value := range picklesList {
			if value.Name == *pickle.Name {
				arr = append(arr, picklesList[key])
				return operations.NewGetPicklesOK().WithPayload(arr)
			}
		}
	}

	return operations.NewGetPicklesOK().WithPayload(picklesList)
}

func GetPickleRandom(pickle operations.GetPickleRandomParams) middleware.Responder {
	var URL string

	// Get Pickles List
	arr := GetPicklesList()

	// Get a Random Index
	rand.Seed(time.Now().UnixNano())
	var index int
	index = rand.Intn(len(arr) - 1)

	URL = "https://raw.githubusercontent.com/cameracode/ricksofpickle/Develop/" + arr[index].Name + ".png"

	response, err := http.Get(URL)
	if err != nil {
		fmt.Println("error")
		srcImage, _ := getFirePickleError("Oops, Error")
		return operations.NewGetPickleNameOK().WithPayload(convertImgToIoCloser(srcImage))
	}
	defer response.Body.Close()

	srcImage, _, err := image.Decode(response.Body)
	if err != nil {
		log.Fatalf("failed to Decode image: %v", err)
	}

	if pickle.Size != nil {
		srcImage = resizeImage(srcImage, *pickle.Size)
	}

	return operations.NewGetPickleNameOK().WithPayload(convertImgToIoCloser(srcImage))
}
