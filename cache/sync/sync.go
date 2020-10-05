package sync

import (
	"encoding/json"
	"fmt"
	"log"
	"pcpartpicker-api/api/endpoints"
	"pcpartpicker-api/api/entities"
	"pcpartpicker-api/api/parse"
	"pcpartpicker-api/cache"
	"time"
)

//I do not think caching CompletedBuilds would be fruitful since there are a lot of variations of the build types cause of Filters and Sorts
// I would keep it to cache when the certain Filter & Sort is used upon.
func Sync() {
	p := parse.Parser{BuildGuides:true}
	_, url := p.ParseToUrl()

	//Ignoring the other return cause at the first init, it's assumed that the cache file is empty (yet).
	//Might change in the future
	guideList, _ := endpoints.GetGuides(url)
	guideB, _ := json.Marshal(guideList)
	ok := cache.Put(url, guideB)
	if !ok {
		log.Println("Error syncing BuildGuide Cache.")
	} else {
		log.Println("Synced BuildGuide")
	}

	for _, cat := range guideList.Categories {
		for _, guide := range cat.Guides {
			syncGuides(guide)
			time.Sleep(time.Minute)
		}
	}



}

func syncGuides(guide entities.Guide) {
	details, _ := endpoints.GetDetails(guide.Path, "")
	detailsB, _ := json.Marshal(details)

	ok := cache.Put(guide.Path, detailsB)
	if !ok {
		log.Println(fmt.Sprintf("Error syncing GuideDetails : %s", guide.Path))
	} else {
		log.Println(fmt.Sprintf("Synced : %s", guide.Path))
	}

	parts, _  := endpoints.GetParts(details.PartsLink)
	partsB, _ := json.Marshal(parts)

	ok = cache.Put(details.PartsLink, partsB)
	if !ok {
		log.Println(fmt.Sprintf("  ↳ Error syncing : %s", details.PartsLink))
	} else {
		log.Println(fmt.Sprintf("  ↳ Synced : %s", details.PartsLink))
	}

}