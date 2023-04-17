package domain

type Alias struct {
	//ID          int
	Name        string `json:"name"`
	Destination string `json:"destination"`
	Group       string `json:"group"`
}
