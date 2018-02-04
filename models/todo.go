package models

type Todo struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}
