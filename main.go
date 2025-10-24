package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type imageFormat string

const (
	FormatJPEG    imageFormat = "jpeg"
	FormatPNG     imageFormat = "png"
	FormatGIF     imageFormat = "gif"
	FormatUnknown imageFormat = "unknown"
)

// toGreyscale converts an image to greyscale using the Luminosity Method.
func toGreyscale(originalImage image.Image) image.Image {
	size := originalImage.Bounds().Size()
	rect := image.Rect(0, 0, size.X, size.Y)

	// Use RGBA to ensure full color range is available for conversion
	modifiedImg := image.NewRGBA(rect)

	for x := 0; x < size.X; x++ {
		for y := 0; y < size.Y; y++ {
			pixel := originalImage.At(x, y)
			// Convert to RGBA to get R, G, B, A
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			red := float64(originalColor.R)
			green := float64(originalColor.G)
			blue := float64(originalColor.B)

			// Simple average method for greyscale
			//grey := uint8(
			//    math.Round((red + green + blue) / 3),
			//)

			// Luminosity method for greyscale
			// Formula: Grey = 0.299*Red + 0.587*Green + 0.114*Blue
			grey := uint8(
				math.Round(0.299*red + 0.587*green + 0.114*blue),
			)

			modifiedColor := color.RGBA{
				R: grey,
				G: grey,
				B: grey,
				A: originalColor.A, // Preserve original alpha channel
			}

			modifiedImg.Set(x, y, modifiedColor)
		}
	}

	return modifiedImg
}

// encodeImage encodes the given image to the specified format
func encodeImage(w *os.File, img image.Image, format imageFormat) error {
	switch format {
	case FormatJPEG:
		return jpeg.Encode(w, img, &jpeg.Options{Quality: 90})
	case FormatPNG:
		return png.Encode(w, img)
	case FormatGIF:
		return gif.Encode(w, img, nil)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <input_image_filename>")
		fmt.Println("Example: go run main.go myimage.jpg")
		os.Exit(1)
	}

	inputFilename := os.Args[1]

	inputFile, err := os.Open(inputFilename)
	if err != nil {
		log.Fatalf("Error opening input file %s: %v", inputFilename, err)
	}
	defer inputFile.Close()

	// Decode the image and determine its format
	img, formatStr, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("Error decoding image %s: %v", inputFilename, err)
	}

	var originalFormat imageFormat
	switch formatStr {
	case "jpeg":
		originalFormat = FormatJPEG
	case "png":
		originalFormat = FormatPNG
	case "gif":
		originalFormat = FormatGIF
	default:
		log.Fatalf("Unsupported input image format: %s. Supported formats are JPEG, PNG, GIF.", formatStr)
	}

	fmt.Printf("Converting %s (format: %s) to greyscale (Luminosity Method)...\n", inputFilename, originalFormat)
	greyImg := toGreyscale(img)
	fmt.Println("Conversion complete.")

	// Construct the output filename
	extension := filepath.Ext(inputFilename)
	baseName := strings.TrimSuffix(inputFilename, extension)
	outputFilename := fmt.Sprintf("%s_greyscale%s", baseName, extension)

	greyscaleOutputFile, err := os.Create(outputFilename)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v", outputFilename, err)
	}
	defer greyscaleOutputFile.Close()

	// Encode using the detected original format
	err = encodeImage(greyscaleOutputFile, greyImg, originalFormat)
	if err != nil {
		log.Fatalf("Error encoding greyscale image to %s (format: %s): %v", outputFilename, originalFormat, err)
	}

	fmt.Printf("Greyscale image saved as %s\n", outputFilename)
}
