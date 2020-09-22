package entities

type Build struct {
	Path      string   `json:"path"`
	Image     string   `json:"image"`
	Author    Author   `json:"author"`
	Title     string   `json:"title"`
	Products  []string `json:"products"`
	Price     string   `json:"price"`
	Followers int      `json:"followers"`
	Comments  int      `json:"comments"`
}

type Author struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
