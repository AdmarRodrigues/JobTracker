package dto

type JobInput struct {
	Cargo   string `json:"cargo" binding:"required"`
	Empresa string `json:"empresa" binding:"required"`
	Status  string `json:"status" binding:"required"`
	Data    string `json:"data" binding:"required"`
}
