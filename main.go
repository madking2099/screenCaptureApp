package main

import (
    "context"
    "fmt"
    "github.com/chromedp/chromedp"
    "github.com/gin-gonic/gin"
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

    // Enhanced CORS with detailed logging
    r.Use(func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        log.Printf("CORS request - Method: %s, Origin: %s, Path: %s", c.Request.Method, origin, c.Request.URL.Path)
        if origin != "" {
            // Allow any origin for testing (restrict in production)
            c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
            c.Writer.Header().Set("Access-Control-Max-Age", "86400") // Cache for 24 hours
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Writer.Header().Set("Vary", "Origin") // Ensure proper caching
        }
        if c.Request.Method == "OPTIONS" {
            log.Println("Handling CORS preflight - 204 No Content")
            c.AbortWithStatus(204)
            return
        }
        c.Next()
    })

    log.Println("Initializing routes")
    r.GET("/", redirectToSwagger)
    log.Println("Serving Swagger UI statically at /swagger/")
    // Dynamically set server host, default to localhost or environment
    serverHost := os.Getenv("SERVER_HOST")
    if serverHost == "" {
        serverHost = "http://localhost:1388" // Default for local testing
        log.Printf("Using default SERVER_HOST: %s", serverHost)
    } else {
        log.Printf("Using SERVER_HOST from environment: %s", serverHost)
    }
    r.Static("/swagger", "./swagger-ui")
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
    log.Println("Received /screenshot request")
    var req ScreenshotRequest
    if err := c.BindJSON(&req); err != nil {
        log.Printf("Failed to bind JSON: %v", err)
        c.JSON(400, gin.H{"detail": "Invalid request"})
        return
    }
    log.Printf("Request bound: %+v", req)

    filename := req.OutputFileName
    if filename == "" {
        filename = fmt.Sprintf("screenshot_%d", os.Getpid())
    }
    if !strings.HasSuffix(filename, ".png") {
        filename += ".png"
    }
    outputFile := filepath.Join("static", filename)
    log.Printf("Output file: %s", outputFile)

    err := captureScreenshot(req.URL, outputFile, req.Headers)
    if err != nil {
        log.Printf("Screenshot failed: %v", err)
        if _, exists := os.Stat(outputFile); exists == nil {
            os.Remove(outputFile)
        }
        c.JSON(500, gin.H{"detail": fmt.Sprintf("Error capturing screenshot: %v", err)})
        return
    }

    log.Println("Screenshot captured successfully")
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
    log.Println("Starting screenshot capture for URL:", url)
    ctx, cancel := chromedp.NewExecAllocator(
        context.Background(),
        append(
            chromedp.DefaultExecAllocatorOptions[:],
            chromedp.Flag("headless", true),
            chromedp.Flag("no-sandbox", true),
            chromedp.Flag("disable-gpu", true),
            chromedp.Flag("ignore-certificate-errors", true),
        )...,
    )
    if ctx == nil {
        log.Println("Failed to create exec allocator")
        return fmt.Errorf("failed to create exec allocator")
    }
    defer cancel()

    ctx, cancel = chromedp.NewContext(ctx)
    if ctx == nil {
        log.Println("Failed to create context")
        return fmt.Errorf("failed to create context")
    }
    defer cancel()

    var buf []byte
    tasks := chromedp.Tasks{
        chromedp.Navigate(url),
        chromedp.WaitVisible("body", chromedp.ByQuery),
        chromedp.FullScreenshot(&buf, 90),
    }
    log.Println("Running chromedp tasks")
    if err := chromedp.Run(ctx, tasks); err != nil {
        log.Printf("Chromedp run failed: %v", err)
        return err
    }

    log.Println("Writing screenshot to file")
    if err := os.WriteFile(outputFile, buf, 0644); err != nil {
        log.Printf("File write failed: %v", err)
        return err
    }

    if len(headers) > 0 {
        log.Printf("Note: Custom headers (%v) are not applied in this version; use URL auth (e.g., https://user:pass@url)", headers)
    }
    return nil
}