package dto

import (
	"fmt"
	"strings"
)

type JobInput struct {
	Cargo   string `json:"cargo" binding:"required"`
	Empresa string `json:"empresa" binding:"required"`
	Status  string `json:"status" binding:"required"`
	Data    string `json:"data" binding:"required"`
}

func (j *JobInput) Validate() error {
	var issues []string

	if strings.TrimSpace(j.Empresa) == "" {
		issues = append(issues, "This field is requested")
	}
	if strings.TrimSpace(j.Cargo) == "" {
		issues = append(issues, "This field is requested")
	}

	if len(issues) > 0 {
		return fmt.Errorf("%s", strings.Join(issues, "; "))
	}
	return nil
}
