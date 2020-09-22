package parse

import (
	"fmt"
	"pcpartpicker-api/api/entities"
)

//Module to parse human-readable and neat structure to dirty PC-PART-PICKER URLs
type Parser struct {
	Region          string                 `json:"region"`
	SystemBuilder   bool                   `json:"system_builder"`
	BuildGuides     bool                   `json:"build_guides"`
	CompletedBuilds completedBuildsOptions `json:"completed_builds"`
}

type completedBuildsOptions struct {
	IsIt          bool                   `json:"is_it"`
	SortOptions   entities.SortOptions   `json:"sort_options"`
	FilterOptions entities.FilterOptions `json:"filter_options"`
}

func (p *Parser) ParseToUrl() (error, string) {
	baseUrl := fmt.Sprintf("https://pcpartpicker.com")

	if p.Region != "" {
		baseUrl = fmt.Sprintf("https://%s.pcpartpicker.com", p.Region)
	}

	if p.SystemBuilder {
		return nil, fmt.Sprint(baseUrl, "/lists/")
	} else if p.BuildGuides {
		return nil, fmt.Sprint(baseUrl, "/guide/")
	} else if p.CompletedBuilds.IsIt {
		baseUrl = fmt.Sprint(baseUrl, "/builds/")

		var sortOption string

		switch {
		case !p.CompletedBuilds.SortOptions.Newest:
			sortOption = "#sort=-recents"
		case p.CompletedBuilds.SortOptions.HighestPriced:
			sortOption = "#sort=price"
		case p.CompletedBuilds.SortOptions.HighestRated:
			sortOption = "#sort=rating"
		default:
			sortOption = "#sort=recents"
		}

		baseUrl = fmt.Sprint(baseUrl, sortOption)

		baseUrl = appendBoolToUrl(baseUrl, "F", p.CompletedBuilds.FilterOptions.Featured)
		baseUrl = appendBoolToUrl(baseUrl, "C", p.CompletedBuilds.FilterOptions.Overclocked)
		baseUrl = appendBoolToUrl(baseUrl, "B", p.CompletedBuilds.FilterOptions.BuildType)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "C", p.CompletedBuilds.FilterOptions.CPUs)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "s", p.CompletedBuilds.FilterOptions.CPUSockets)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "h", p.CompletedBuilds.FilterOptions.CPUCoolers)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "g", p.CompletedBuilds.FilterOptions.GPUs)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "G", p.CompletedBuilds.FilterOptions.SLI)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "e", p.CompletedBuilds.FilterOptions.Case)
		baseUrl = appendBasicOptionsToUrl(baseUrl, "E", p.CompletedBuilds.FilterOptions.CaseType)

	}
	return nil, baseUrl
}

func appendBoolToUrl(url string, param string, value *bool) string {
	x := 0
	if value == nil {
		return url
	}

	if *value {
		x = 1
	}
	return url + fmt.Sprintf("&%s=%d", param, x)
}

func appendBasicOptionsToUrl(url, param string, values []entities.BasicOptions) string {
	firstTime := true

	for _, value := range values {
		if value.Selected && firstTime {
			url = fmt.Sprint(url, fmt.Sprintf("&c=%s", value.Path))
			firstTime = false
		} else {
			url = fmt.Sprint(url, fmt.Sprintf(",%s", value.Path))
		}
	}

	return url
}
