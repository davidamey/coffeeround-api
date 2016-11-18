package models

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	TotalDay  int    `json:"totalDay"`
	TotalWeek int    `json:"totalWeek"`
	Total     int    `json:"total"`
}
