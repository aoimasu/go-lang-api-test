package models

type Tag struct {
	Tag     			string 			`json:"tag"`
	Count   			int 				`json:"count"`
	Articles     	[]string 		`json:"articles"`
	RelatedTags   []string 		`json:"related_tags"`
}