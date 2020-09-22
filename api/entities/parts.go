package entities

type Parts struct {
	Compatibility bool         `json:"compatibility"`
	Wattage       string       `json:"wattage"`
	CPU           partsDetails `json:"cpu"`
	Motherboard   partsDetails `json:"motherboard"`
	Memory        partsDetails `json:"memory"`
	Storage       partsDetails `json:"storage"`
	VideoCard     partsDetails `json:"video_card"`
	Case          partsDetails `json:"case"`
	PSU           partsDetails `json:"psu"`
}

type partsDetails struct {
	Title string `json:"title"`
	Image string `json:"image"`
	Price string `json:"price"`
	Where string `json:"where"`
}
