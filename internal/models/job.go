package models

type Jobs struct {
	Id      int    `json:"id"`
	Cargo   string `json:"cargo"`
	Empresa string `json:"empresa"`
	Status  string `json:"status"`
	Data    string `json:"data"`
}
