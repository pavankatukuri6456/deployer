package main

import (
	"context"
	"deployer/docker"
	"deployer/model"
	"net/http"

	"github.com/gin-gonic/gin"
	// Replace with your actual module path
)

// @Summary Trigger a pipeline
// @Description Trigger a Tekton pipeline with provided Git repo, branch, app name and instance.
// @Accept  json
// @Produce  json
// @Param   request body DeployRequest true "Deployment configuration"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /deploy [post]
func main() {

	r := gin.Default()

	r.POST("/deploy", func(c *gin.Context) {
		var req model.DeployRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := context.TODO()
		prName, imageURL, err := docker.TriggerPipeline(ctx, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"pipelineRun": prName,
			"image":       imageURL,
			"status":      "Pipeline triggered",
		})
	})

	r.GET("/deploy", Getstatus)

	r.Run(":8080")

}

func Getstatus(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}
