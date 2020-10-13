package documents

//import "encoding/json"

type Document struct {
	ID		string `json:"id"`
	Title	string `json:"title"`
	Author	string `json:"author"`
	Summary	string `json:"summary"`
}
