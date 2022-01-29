package asciiimageservice

import (
	"errors"
	"image"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/qeesung/image2ascii/convert"
)

const IMAGES_DIR = "images"

type AsciiImage struct {
	Id, Data string
}

var ImageNotFound = errors.New("Image not found")

type AsciiImageService interface {
	SaveAsAscii(io.Reader) (string, error)
	GetById(string) (*AsciiImage, error)
	ListAllIds() ([]string, error)
}

type AsciiImageFileService struct {
	imagetoasciiconverter *convert.ImageConverter
}

func New() *AsciiImageFileService {
	initImagesDir()
	return &AsciiImageFileService{
		imagetoasciiconverter: convert.NewImageConverter(),
	}
}

func initImagesDir() {
	_, err := os.Stat(IMAGES_DIR)
	if err == nil {
		return
	} else {
		if errors.Is(err, os.ErrNotExist) {
			err = os.Mkdir("images", os.FileMode(0755))
			if err != nil {
				log.Fatal("Error creating images directory", err)
			}
		} else {
			log.Fatal("Error checking if images directory exists", err)
		}
	}
}

func (a *AsciiImageFileService) SaveAsAscii(data io.Reader) (string, error) {
	imageId := uuid.New().String()

	img, _, err := image.Decode(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	asciiImage := a.imagetoasciiconverter.Image2ASCIIString(img, &convert.DefaultOptions)

	file, err := os.Create(IMAGES_DIR + "/" + imageId + ".txt")
	defer file.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	_, err = file.WriteString(asciiImage)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return imageId, nil
}

func (a *AsciiImageFileService) GetById(imageId string) (*AsciiImage, error) {
	files, err := ioutil.ReadDir(IMAGES_DIR)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, file := range files {
		if strings.TrimRight(file.Name(), ".txt") == imageId {
			image, err := ioutil.ReadFile(IMAGES_DIR + "/" + file.Name())
			if err != nil {
				log.Println(err)
				return nil, err
			}

			return &AsciiImage{Id: imageId, Data: string(image)}, nil
		}
	}

	return nil, ImageNotFound
}

func (a *AsciiImageFileService) ListAllIds() ([]string, error) {
	files, err := ioutil.ReadDir(IMAGES_DIR)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	imageIds := make([]string, len(files))
	for i, file := range files {
		imageIds[i] = strings.TrimRight(file.Name(), ".txt")
	}

	return imageIds, nil
}
