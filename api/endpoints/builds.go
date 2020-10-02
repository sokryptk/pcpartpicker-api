package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"log"
	"net/http"
	"pcpartpicker-api/api/entities"
	"pcpartpicker-api/api/parse"
	"pcpartpicker-api/cache"
	"pcpartpicker-api/scraper"
	"strconv"
	"sync"
)

func GetCompletedBuilds(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rBody, _ := ioutil.ReadAll(r.Body)

	region := r.Header.Get("region")

	// Options are 0, 1, 2
	// 0 - Newest
	// 1 - Highest Rated
	// 2 - Highest Priced
	// Default is Newest
	sorted, _ := strconv.Atoi(r.Header.Get("sort"))

	var filterOptions entities.FilterOptions

	_ = json.Unmarshal(rBody, &filterOptions)

	var sortOptions entities.SortOptions

	switch sorted {
	case 0:
		sortOptions.Newest = true
	case 1:
		sortOptions.HighestRated = true
	case 2:
		sortOptions.HighestPriced = true
	default:
		sortOptions.Newest = true
	}

	p  := parse.Parser{
		Region:          region,
		CompletedBuilds: parse.CompletedBuildsOptions{
			IsIt:true,
			FilterOptions:filterOptions,
			SortOptions: sortOptions,
		},
	}

	_, url := p.ParseToUrl()

	if data, success := cache.RetrieveCache(url); success {
		var db []entities.Build

		_ = json.Unmarshal(data, &db)

		_ = json.NewEncoder(w).Encode(db)

		return
	}


	if _, err := scraper.Instance.ExecuteScript(fmt.Sprintf("window.open('%s');", url), nil); err != nil {
		log.Println(err)
	}

	windows, _ := scraper.Instance.WindowHandles()
	_ = scraper.Instance.SwitchWindow(windows[len(windows)-1])

	handle, _ := scraper.Instance.CurrentWindowHandle()
	defer scraper.Instance.SwitchWindow(windows[0])
	defer scraper.Instance.CloseWindow(handle)

	_ = scraper.Instance.Wait(func(wd selenium.WebDriver) (b bool, err error) {
		e, _ := wd.FindElements(selenium.ByCSSSelector, ".logGroup__card")
		if len(e) > 0 {
			return true, nil
		}

		return false, nil
	})

	var builds []entities.Build

	buildCards, _ := scraper.Instance.FindElements(selenium.ByCSSSelector, ".logGroup__card")

	var wg sync.WaitGroup

	for _, card := range buildCards {
		var build entities.Build

		wg.Add(1)
		go appendEntitiesToBuild(card, &build, &wg)
		wg.Wait()

		builds = append(builds, build)
	}

	_ = json.NewEncoder(w).Encode(builds)

	b, _ := json.Marshal(builds)
	cache.Put(url, b)

}

func appendEntitiesToBuild(card selenium.WebElement, build *entities.Build, wg *sync.WaitGroup) {
	defer wg.Done()
	path, _ := card.FindElement(selenium.ByCSSSelector, "a")
	build.Path, _ = path.GetAttribute("href")

	price, _ := card.FindElement(selenium.ByCSSSelector, ".log__price")
	build.Price, _  = price.Text()

	comments, _ := card.FindElement(selenium.ByCSSSelector, ".log__link--comments")
	commentTxt, _  := comments.Text()
	build.Comments, _ = strconv.Atoi(commentTxt)

	followers, _ := card.FindElement(selenium.ByCSSSelector, ".log__link--followers")
	followersTxt, _ := followers.Text()
	build.Followers, _ = strconv.Atoi(followersTxt)

	author, _ := card.FindElements(selenium.ByCSSSelector, ".log__author a")
	//The first object is the one having the user avatar. Hence, un-needed and thus skipped.
	build.Author.Path, _ = author[1].GetAttribute("href")
	build.Author.Name, _ = author[1].Text()

	title, _ := card.FindElement(selenium.ByCSSSelector, ".log__title a")
	build.Title, _ = title.Text()

	products, _ := card.FindElements(selenium.ByCSSSelector, "build__specs")
	for _, product := range products {
		productTxt, _ := product.Text()
		build.Products = append(build.Products, productTxt)
	}

	//I think for an app, images here, are un-necessary. Since this would be a list and images would make it far cluttered.
	// Not removing from the entities for now, will decide in the near future.

}