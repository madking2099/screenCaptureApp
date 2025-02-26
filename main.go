package main

import (
    "context"
    "fmt"
    "github.com/chromedp/chromedp"
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
    "log"
    "os"
    "path/filepath"
    "strings"
)

// @title Webpage Screenshot Service
// @version 1.0.0
// @description Capture screenshots of webpages as images and serve them.
// @license.name MIT
// @license.url https://github.com/madking2099/screenCaptureApp/blob/master/LICENSE
// @host localhost:1388
// @BasePath /

type ScreenshotRequest struct {
    URL           string            `json:"url" binding:"required" example:"https://example.com"`
    Headers       map[string]string `json:"headers" example:"{\"User-Agent\": \"MyBot/1.0\"}"`
    OutputFileName string            `json:"output_filename" example:"screenshot"`
}

func main() {
    r := gin.Default()
    r.Static("/static", "./static")
    if err := os.MkdirAll("static", 0755); err != nil {
        log.Fatal(err)
    }

    r.GET("/", func(c *gin.Context) {
        c.Redirect(302, "/docs")
    })

    r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // @Summary Check service health
    // @Description Returns the health status of the service
    // @Produce json
    // @Success 200 {object} map[string]string
    // @Router /health [get]
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

    // @Summary Capture a webpage screenshot
    // @Description Takes a URL and returns a screenshot file URL
    // @Accept json
    // @Produce json
    // @Param request body ScreenshotRequest true "Screenshot request payload"
    // @Success 200 {object} map[string]string
    // @Failure 400 {object} map[string]string
    // @Failure 500 {object} map[string]string
    // @Router /screenshot [post]
    r.POST("/screenshot", func(c *gin.Context) {
        var req ScreenshotRequest
        if err := c.BindJSON(&req); err != nil {
            c.JSON(400, gin.H{"detail": "Invalid request"})
            return
        }

        filename := req.OutputFileName
        if filename == "" {
            filename = fmt.Sprintf("screenshot_%d", os.Getpid())
        }
        if !strings.HasSuffix(filename, ".png") {
            filename += ".png"
        }
        outputFile := filepath.Join("static", filename)

        err := captureScreenshot(req.URL, outputFile, req.Headers)
        if err != nil {
            if _, exists := os.Stat(outputFile); exists == nil {
                os.Remove(outputFile)
            }
            c.JSON(500, gin.H{"detail": fmt.Sprintf("Error capturing screenshot: %v", err)})
            return
        }

        c