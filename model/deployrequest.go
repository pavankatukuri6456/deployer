package model

type DeployRequest struct {
	RepoURL  string `json:"repo_url" binding:"required"`
	Branch   string `json:"branch" binding:"required"`
	AppName  string `json:"app_name" binding:"required"`
	Instance string `json:"instance" binding:"required"`
}
