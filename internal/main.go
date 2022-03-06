package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"pickle-api/pkg/swagger/server/models"
	"pickle-api/pkg/swagger/server/restapi"
	"pickle-api/pkg/swagger/server/restapi/operations"

	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/google/go-github/v43/github"
)

func main() {
	// Initialize Swagger
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewPickleAPIAPI(swaggerSpec)
	// Use SwaggerUI instead of reDoc on /docs
	api.UseSwaggerUI()

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

	api.GetPicklesHandler = operations.GetPicklesHandlerFunc(GetPickles)

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
		srcImage, _ := getBladedPickleError("Oops, Error")
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
		srcImage, _ := getBladedPickleError("Oops, Error")
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

/*
Display Bladed Pickle Rick with a message (error)
*/
func getBladedPickleError(message string) (image.Image, error) {
	// open local file
	file, err := os.Open("./assets/bladed_picklerick.png")
	if err != nil {
		log.Fatalf("failed to Open bladed_picklerick image: %v", err)
	}

	srcImage, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to Deecode image: %v", err)
		return srcImage, err
	}
	// Add text on Pickle Rick
	srcImage, err = TextOnPickleRick(srcImage, "Oops, Error! Pickle Rick has blades, dawg!")

	// Resize Image
	srcImage = resizeImage(srcImage, "medium")
	if err != nil {
		log.Fatalf("failed to put Text on Pickle Rick: %v", err)
		return srcImage, err
	}

	return srcImage, nil
}

/*
Get List of Pickle Ricks from cameracode repo
*/
func GetPicklesList() []*models.Pickle {
	client := github.NewClient(nil)

	// list public repos for org "github"
	ctx := context.Background()
	// list all repos for the authenticated user
	_, directoryContent, _, err := client.Repositories.GetContents(ctx, "cameracode", "ricksofpickle", "/", nil)
	if err != nil {
		fmt.Println(err)
	}

	var arr []*models.Pickle

	for _, c := range directoryContent {
		if *c.Name == ".gtiignore" || *c.Name == "README.md" {
			continue
		}

		var name string = strings.Split(*c.Name, ".")[0]
		arr = append(arr, &models.Pickle{name, *c.Path, *c.DownloadURL})
	}
	return arr
}

/*
Resize Image
*/
func resizeImage(srcImage image.Image, size string) image.Image {
	var height int
	switch size {
	case "x-small":
		height = 50
	case "small":
		height = 100
	case "medium":
		height = 300
	default:
		// Lulzbotz
		height = 1000
	}

	// Resize the cropped image to width = 200px preserving the aspect ratio
	srcImage = imaging.Resize(srcImage, 0, height, imaging.Lanczos)

	return srcImage
}

/*
Convert Image to io.close (for reply format)
*/
func convertImgToIoCloser(srcImage image.Image) io.ReadCloser {
	encoded := &bytes.Buffer{}
	png.Encode(encoded, srcImage)

	return ioutil.NopCloser(encoded)
}

/*
Add text on Image
*/
func TextOnPickleRick(bgImage image.Image, text string) (image.Image, error) {
	imgWidth := bgImage.Bounds().Dx()
	imgHeight := bgImage.Bounds().Dy()

	dc := gg.NewContext(imgWidth, imgHeight)
	dc.DrawImage(bgImage, 0, 0)

	if err := dc.LoadFontFace("assets/FireSans-Light.ttf", 50); err != nil {
		return nil, err
	}

	x := float64((imgWidth / 2))
	y := float64((imgHeight / 12))
	maxWidth := float64(imgWidth) - 60.0
	dc.SetColor(color.Black)
	dc.DrawStringWrapped(text, x, y, 0.5, 0.5, maxWidth, 1.5, gg.AlignRight)

	return dc.Image(), nil
}
