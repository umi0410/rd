package domain

type Alias struct {
	//ID          int
	Group       string `json:"group"`
	Name        string `json:"name"`
	Destination string `json:"destination"`
}
