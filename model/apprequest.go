package model

type AppRequest struct {
	AppName   string `json:"appName"`
	Namespace string `json:"namespace"`
	Image     string `json:"image"`
}
