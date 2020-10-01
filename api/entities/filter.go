package entities

type FilterOptions struct {
	Price       priceOptions     `json:"price"`
	Published   publishedOptions `json:"published"`
	Featured    *bool            `json:"featured"`
	BuildType   *bool `json:"build_type"`
	CPUs        []BasicOptions   `json:"cpus"`
	Overclocked *bool            `json:"overclocked"`
	CPUSockets  []BasicOptions   `json:"cpu_sockets"`
	CPUCoolers  []BasicOptions   `json:"cpu_coolers"`
	GPUs        []BasicOptions   `json:"gpus"`
	SLI         []BasicOptions   `json:"sli"`
	Case        []BasicOptions   `json:"case"`
	CaseType    []BasicOptions   `json:"case_type"`
}

type SortOptions struct {
	Newest        bool `json:"newest"`
	HighestRated  bool `json:"highest_rated"`
	HighestPriced bool `json:"highest_priced"`
}

type priceOptions struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type publishedOptions struct {
	All  bool `json:"all"`
	OneD bool `json:"one_d"`
	OneW bool `json:"one_w"`
	OneM bool `json:"one_m"`
	OneY bool `json:"one_y"`
}

type BasicOptions struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Selected bool   `json:"selected"`
}
