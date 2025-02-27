package main

import (
    "context"
    "fmt"
    "github.com/chromedp/chromedp"
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
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

    log.Println("Initializing routes")
    r.GET("/", redirectToSwagger)
    // Serve Swagger UI at /swagger/ with wildcard
    log.Println("Registering Swagger UI at /swagger/*any with JSON at /api-docs/swagger.json")
    swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api-docs/swagger.json"))
    r.GET("/swagger/*any", func(c *gin.Context) {
        log.Printf("Handling Swagger request for path: %s", c.Request.URL.Path)
        swaggerHandler(c)
    })
    r.StaticFile("/api-docs/swagger.json", "./docs/swagger.json")
    r.GET("/health", getHealth)
    r.POST("/screenshot", postScreenshot)
    r.DELETE("/static/:filename", deleteScreenshot)

    r.Run(":8000")
}

func redirectToSwagger(c *gin.Context) {
    c.Redirect(302, "/swagger/")
}

// @Summary Check service health
// @Description Returns the health status of the service
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func getHealth(c *gin.Context) {
    c.JSON(200, gin.H{"status": "healthy"})
}

// @Summary Capture a webpage screenshot
// @Description Takes a URL and returns a screenshot file URL. Headers beyond basic auth in URL are not supported yet.
// @Accept json
// @Produce json
// @Param request body ScreenshotRequest true "Screenshot request payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /screenshot [post]
func postScreenshot(c *gin.Context) {
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

    c.JSON(200, gin.H{"file_url": fmt.Sprintf("/static/%s", filename)})
}

// @Summary Delete a screenshot file
// @Description Removes a screenshot file from the static directory
// @Produce json
// @Param filename path string true "Filename to delete"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /static/{filename} [delete]
func deleteScreenshot(c *gin.Context) {
    filename := c.Param("filename")
    filePath := filepath.Join("static", filename)
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(404, gin.H{"detail": "File not found"})
        return
    }
    if err := os.Remove(filePath); err != nil {
        c.JSON(500, gin.H{"detail": fmt.Sprintf("Error deleting file: %v", err)})
        return
    }
    c.JSON(200, gin.H{"message": fmt.Sprintf("File %s deleted", filename)})
}

func captureScreenshot(url, outputFile string, headers map[string]string) error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("ignore-certificate-errors", true))
    ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
    defer cancel()

    var buf []byte
    tasks := chromedp.Tasks{
        chromedp.Navigate(url),
        chromedp.WaitVisible("body", chromedp.ByQuery),
        chromedp.FullScreenshot(&buf, 90),
    }

    if err := chromedp.Run(ctx, tasks); err != nil {
        return err
    }
    if len(headers) > 0 {
        log.Printf("Note: Custom headers (%v) are not applied in this version; use URL auth (e.g., https://user:pass@url)", headers)
    }
    return os.WriteFile(outputFile, buf, 0644)
}