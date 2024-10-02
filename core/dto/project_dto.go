package dto

type ProjectDTO struct {
	Title       string `json:"title" required:"true"`
	Description string `json:"description" required:"false"`
	Done        bool   `json:"done" required:"false"`
}

type UpdateProjectDTO struct {
	Id          int64  `json:"id" required:"true"`
	Title       string `json:"title" required:"false"`
	Description string `json:"description" required:"false"`
	Done        bool   `json:"done" required:"false"`
}
