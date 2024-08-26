package main

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "github.com/unidoc/unioffice/presentation"
    "io/ioutil"
    "net/http"
    "path/filepath"
)

// SlideData represents the structure of a slide's content
type SlideData struct {
    SlideNumber   int             `json:"slide_number"`             // The number of the slide in the presentation
    Title         string          `json:"title,omitempty"`          // The main title of the slide
    Subheading    string          `json:"subheading,omitempty"`     // The optional subheading of the slide
    Paragraphs    []string        `json:"paragraphs,omitempty"`     // Any text content on the slide
    Images        []ImageData     `json:"images,omitempty"`         // List of images on the slide
    SpeakerNotes  string          `json:"speaker_notes,omitempty"`  // Speaker notes for the slide
}

// ImageData represents the structure of an image within a slide
type ImageData struct {
    Src       string `json:"src,omitempty"`        // The file name or path of the image
    Alt       string `json:"alt,omitempty"`        // Alt text for the image (for accessibility)
    Width     int    `json:"width,omitempty"`      // The width of the image
    Height    int    `json:"height,omitempty"`     // The height of the image
    Base64    string `json:"base64,omitempty"`     // The Base64-encoded image data
}

// encodeImageToBase64 reads an image file and encodes it as a Base64 string
func encodeImageToBase64(imagePath string) (string, error) {
    // Read the image file into memory
    imgData, err := ioutil.ReadFile(imagePath)
    if err != nil {
        return "", err
    }
    // Encode the image data as a Base64 string
    return base64.StdEncoding.EncodeToString(imgData), nil
}

// parsePPTX parses the PowerPoint file and extracts slide content
func parsePPTX(filePath string) ([]SlideData, error) {
    // Open the PowerPoint presentation
    ppt, err := presentation.Open(filePath)
    if err != nil {
        return nil, err
    }
    defer ppt.Close() // Ensure the file is closed when done

    var slidesContent []SlideData

    // Iterate through each slide in the presentation
    for i, slide := range ppt.Slides() {
        slideData := SlideData{
            SlideNumber: i + 1, // Slide numbers are 1-based
        }

        // Extract text and classify it as Title, Subheading, or Paragraphs
        for _, shape := range slide.Shapes() {
            if shape.HasText() {
                text, _ := shape.Text()
                // Simple heuristic to classify short text as titles and longer text as paragraphs
                if slideData.Title == "" && len(text) < 50 {
                    slideData.Title = text
                } else {
                    slideData.Paragraphs = append(slideData.Paragraphs, text)
                }
            }
        }

        // Extract and encode images as Base64
        for _, pic := range slide.Pictures() {
            imgData, _ := pic.Data()
            imgFileName := filepath.Base(pic.FileName())
            // Save the image file locally
            ioutil.WriteFile(imgFileName, imgData, 0644)

            // Encode the image as Base64
            base64Data, err := encodeImageToBase64(imgFileName)
            if err != nil {
                fmt.Println("Error encoding image:", err)
                continue
            }

            // Add the image data to the slide content
            slideData.Images = append(slideData.Images, ImageData{
                Src:    imgFileName,
                Base64: base64Data,
            })
        }

        // Add the slide data to the list of slides
        slidesContent = append(slidesContent, slideData)
    }
    return slidesContent, nil
}

func main() {
    // Initialize the Gin router
    r := gin.Default()

    // Define the POST /parse endpoint
    r.POST("/parse", func(c *gin.Context) {
        // Retrieve the uploaded file
        file, _ := c.FormFile("file")
        // Save the uploaded file locally
        err := c.SaveUploadedFile(file, file.Filename)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to save file"})
            return
        }

        // Parse the PowerPoint file to extract slide content
        slidesContent, err := parsePPTX(file.Filename)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse PowerPoint"})
            return
        }

        // Respond with the slide content as JSON
        c.JSON(http.StatusOK, slidesContent)
    })

    // Start the HTTP server on port 8080
    r.Run(":8080")
}
