package dto

import (
	"errors"
)

type ProjectDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type UpdateProjectDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Done        *bool   `json:"done"`
}

func (up *UpdateProjectDTO) Validate() error {
	if up.Title == nil && up.Description == nil && up.Done == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
