package endpoints

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"log"
	"net/http"
	"pcpartpicker-api/api/entities"
	"pcpartpicker-api/api/parse"
	"pcpartpicker-api/cache"
	"pcpartpicker-api/scraper"
	"strconv"
	"time"
)

func GetBuildGuides(w http.ResponseWriter, r *http.Request) {
	region := r.Header.Get("region")

	options := parse.Parser{BuildGuides: true, Region: region}
	_, url := options.ParseToUrl()

	guidesList, cached := GetGuides(url)

	_ = json.NewEncoder(w).Encode(guidesList)

	if !cached {
		b, _ := json.Marshal(guidesList)
		cache.Put(url, b)
	}
}

//Returns true if value is retrieved from Cache, else false
func GetGuides(url string) (entities.GuideList, bool) {
	if data, success := cache.RetrieveCache(url); success {
		var db entities.GuideList

		_ = json.Unmarshal(data, &db)

		return db, true
	}

	if _, err := scraper.Instance.ExecuteScript(fmt.Sprintf("window.open('%s');", url), nil); err != nil {
		log.Println(err)
	}

	windows, _ := scraper.Instance.WindowHandles()
	_ = scraper.Instance.SwitchWindow(windows[len(windows)-1])

	handle, _ := scraper.Instance.CurrentWindowHandle()
	defer scraper.Instance.SwitchWindow(windows[0])
	defer scraper.Instance.CloseWindow(handle)

	var guidesList entities.GuideList

	err := scraper.Instance.WaitWithTimeout(func(wd selenium.WebDriver) (b bool, err error) {
		main, _  := wd.FindElements(selenium.ByCSSSelector, ".guideGroup")

		if len(main) > 0 {
			return true, nil
		}

		return false, nil
	}, time.Minute)

	if err != nil {
		log.Println(err)
	}

	guides, _ := scraper.Instance.FindElements(selenium.ByCSSSelector, ".main-content .block")

	for _, subGuides := range guides {
		category := struct {
			Title  string           `json:"title"`
			Guides []entities.Guide `json:"guides"`
		}{}

		categoryTitle, _ := subGuides.FindElement(selenium.ByCSSSelector, "h2")
		category.Title, _ = categoryTitle.Text()
		cards, _ := subGuides.FindElements(selenium.ByCSSSelector, ".guideGroup.guideGroup__card")

		for _, card := range cards {
			guide := entities.Guide{}
			guideTitle, _ := card.FindElement(selenium.ByCSSSelector, ".guide__title")
			guide.Title, _ = guideTitle.Text()

			path, _ := card.FindElement(selenium.ByCSSSelector, ".guideGroup__target")
			guide.Path, _ = path.GetAttribute("href")

			guideProducts, _ := card.FindElements(selenium.ByCSSSelector, ".guide__keyProducts li")

			for _, prod := range guideProducts {
				//I'm sick of the variable names, lol.
				t, _ := prod.Text()

				guide.Products = append(guide.Products, t)
			}

			price, _ := card.FindElement(selenium.ByCSSSelector, ".guide__price")
			guide.Price, _ = price.Text()

			comments, _ := card.FindElement(selenium.ByCSSSelector, ".guide__link--comments")
			cP, _ := comments.Text()
			guide.Comments, _ = strconv.Atoi(cP)

			images, _ := card.FindElements(selenium.ByCSSSelector, ".guide__images img")
			for _, image := range images {
				gI, _ := image.GetAttribute("src")
				guide.Images = append(guide.Images, gI)
			}
			category.Guides = append(category.Guides, guide)
		}

		guidesList.Categories = append(guidesList.Categories, category)
	}

	return guidesList, false
}
