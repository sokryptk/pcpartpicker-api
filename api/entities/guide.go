package entities

type GuideList struct {
	Categories []struct {
		Title  string  `json:"title"`
		Guides []Guide `json:"guides"`
	} `json:"categories"`
}

type Guide struct {
	Path     string   `json:"path"`
	Title    string   `json:"title"`
	Products []string `json:"products"`
	Price    string   `json:"price"`
	Comments int      `json:"comments"`
	Images   []string `json:"images"`
}

type GuideDetails struct {
	Images           []string         `json:"images"`
	Votes            int              `json:"votes"`
	NumberOfComments int              `json:"comments"`
	PartsLink        string           `json:"parts_link"`
	Description      guideDescription `json:"description"`
	Comments         guideComments    `json:"comments"`
}

type guideDescription struct {
	CPU          string `json:"cpu"`
	Motherboard  string `json:"motherboard"`
	Memory       string `json:"memory"`
	Storage      string `json:"storage"`
	GPU          string `json:"gpu"`
	Case         string `json:"case"`
	PSU          string `json:"psu"`
}

type guideComments struct {
}
