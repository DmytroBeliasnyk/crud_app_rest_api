package dto

type ProjectDTO struct {
	Title       string `json:"title" required:"true"`
	Description string `json:"description" required:"false"`
	Done        bool   `json:"done" required:"false"`
}

type UpdateProjectDTO struct {
	Title       string `json:"title" required:"false"`
	Description string `json:"description" required:"false"`
	Done        bool   `json:"done" required:"false"`
}
