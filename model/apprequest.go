package model

type AppRequest struct {
	AppName   string `json:"appName" binding:"required"`
	Namespace string `json:"namespace" binding:"required"`
	Image     string `json:"image" binding:"required"`
}
