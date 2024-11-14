package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"ykkalexx.com/wowraidhelper/internal/repository"
	"ykkalexx.com/wowraidhelper/internal/service"
	"ykkalexx.com/wowraidhelper/pkg/wowapi"
)

func main() {
	// init wow client
	wowClient := wowapi.NewClient(
		os.Getenv("WOW_CLIENT_ID"),
		os.Getenv("WOW_CLIENT_SECRET"),
	)

	// init repo
    repo, err := repository.NewRepository(
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"), 
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )
    if err != nil {
        log.Fatal("failed to init repo:", err)
    }
    defer repo.Close()
	
	// init service
	raidService := service.NewRaidService(wowClient, repo)

	//setup router
	r := gin.Default()

	// endpoint
	r.POST("/analyze-raid", func (c *gin.Context) {
		var request service.RaidAnalysisRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
            return
		}

		analysis, err := raidService.AnalyzeRaid(request)
        if err != nil {
            c.JSON(500, gin.H{"error": err.Error()})
            return
        }

		c.JSON(200, analysis)
	})

	r.Run(":8080")
}