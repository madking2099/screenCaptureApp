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

type ScreenshotRequest struct {
    URL           string            `json:"url" binding:"required"`
    Headers       map[string]string `json:"headers"`
    OutputFileName string            `json:"output_filename"`
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

    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "healthy"})
    })

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

        c.JSON(200, gin.H{"file_url": fmt.Sprintf("/static/%s", filename)})
    })

    r.DELETE("/static/:filename", func(c *gin.Context) {
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
    })

    r.Run(":8000")
}

func captureScreenshot(url, outputFile string, headers map[string]string) error {
    ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    // Ignore SSL errors for self-signed certs
    opts := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("ignore-certificate-errors", true))
    ctx, cancel = chromedp.NewExecAllocator(ctx, opts...)
    defer cancel()

    var buf []byte
    tasks := chromedp.Tasks{
        chromedp.Navigate(url),
        chromedp.WaitVisible("body", chromedp.ByQuery),
        chromedp.FullScreenshot(&buf, 90),
    }
    if headers != nil {
        tasks = append(chromedp.Tasks{chromedp.ActionFunc(func(ctx context.Context) error {
            return chromedp.Run(ctx, chromedp.SetExtraHTTPHeaders(headers))
        })}, tasks...)
    }

    if err := chromedp.Run(ctx, tasks); err != nil {
        return err
    }
    return os.WriteFile(outputFile, buf, 0644)
}