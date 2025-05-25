package server

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine) {
	r.POST("/deploy", DeployApp)
}
