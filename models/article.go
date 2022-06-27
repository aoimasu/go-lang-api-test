package models

type Article struct {
	Id      int64 		`json:"id,string,omitempty"`
	Title   string 		`json:"title,omitempty"`
	Date    string 		`json:"date,omitempty"`
	Body    string 		`json:"body,omitempty"`
	Tags 	  []string 	`json:"tags,omitempty"`
}