package entity

type Project struct {
	Id          int    `json:"id" db:"id" required:"false"`
	Title       string `json:"title" db:"title" required:"true"`
	Description string `json:"description" db:"description" required:"false"`
	Done        bool   `json:"done" db:"done" required:"false"`
}
