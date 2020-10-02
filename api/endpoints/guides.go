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
	"strings"
	"time"
)

func GetBuildGuides(w http.ResponseWriter, r *http.Request) {
	region := r.Header.Get("region")

	options := parse.Parser{BuildGuides: true, Region: region}
	_, url := options.ParseToUrl()

	if data, success := cache.RetrieveCache(url); success {
		var db entities.GuideList

		_ = json.Unmarshal(data, &db)

		_ = json.NewEncoder(w).Encode(db)

		return
	}

	if err := scraper.Instance.Get(url); err != nil {
		log.Println(err)
	}

	handle, _ := scraper.Instance.CurrentWindowHandle()
	defer scraper.Instance.CloseWindow(handle)

	var guidesList entities.GuideList

	err := scraper.Instance.WaitWithTimeout(func(wd selenium.WebDriver) (b bool, err error) {
		main, _  := wd.FindElements(selenium.ByCSSSelector, ".main-content .block")

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

	_ = json.NewEncoder(w).Encode(guidesList)

	if guidesList.Categories != nil {
		b, _ := json.Marshal(guidesList)
		cache.Put(url, b)
	}
}

func GetGuideDetails(w http.ResponseWriter, r *http.Request) {
	path := r.Header.Get("path")
	region := r.Header.Get("region")

	if data, success := cache.RetrieveCache(path); success {
		var db entities.GuideDetails

		_ = json.Unmarshal(data, &db)

		_ = json.NewEncoder(w).Encode(db)

		return
	}

	err := scraper.Instance.Get(path)
	if err != nil {
		log.Println(err)
	}

	handle, _ := scraper.Instance.CurrentWindowHandle()
	defer scraper.Instance.CloseWindow(handle)

	var guideDetails entities.GuideDetails

	err = scraper.Instance.WaitWithTimeout(func(wd selenium.WebDriver) (b bool, err error) {
		i, _ := wd.FindElements(selenium.ByCSSSelector, ".actionBoxGroup, .description")
		if len(i) > 1 {
			return true, nil
		}

		return false, nil
	}, time.Minute)

	if err != nil {
		log.Println(err)
	}

	images, _ := scraper.Instance.FindElements(selenium.ByCSSSelector, ".gallery__image")
	for _, image := range images {
		src, _ := image.FindElement(selenium.ByCSSSelector, "img")

		imgSrc, _ := src.GetAttribute("src")
		guideDetails.Images = append(guideDetails.Images, imgSrc)
	}

	vote, _ := scraper.Instance.FindElement(selenium.ByCSSSelector, ".actionBox__vote span")
	vText, _ := vote.Text()
	guideDetails.Votes, _ = strconv.Atoi(vText)

	link, _ := scraper.Instance.FindElement(selenium.ByCSSSelector, ".subTitle__form a")
	pLink, _ := link.GetAttribute("onclick")
	pUrl := parse.Parser{Region:region}
	_, url := pUrl.ParseToUrl()

	guideDetails.PartsLink = fmt.Sprint(url, fmt.Sprintf("/list/%s", strings.Split(pLink, "'")[1]))
	des, _ := scraper.Instance.FindElement(selenium.ByCSSSelector, ".description")

	desElements, _ := des.FindElements(selenium.ByCSSSelector, "h2, p")
	isPart := 0
	for _, e := range desElements {
		tag, _ := e.TagName()
		if tag == "h2" {
			isPart = isPart + 1
		}

		text, _ := e.Text()

		if tag != "h2" {
			switch isPart {
			case 1:
				continue
			case 2:
				guideDetails.Description.CPU = fmt.Sprint(guideDetails.Description.CPU, text)
			case 3:
				guideDetails.Description.Motherboard = fmt.Sprint(guideDetails.Description.Motherboard, text)
			case 4:
				guideDetails.Description.Memory = fmt.Sprint(guideDetails.Description.Memory, text)
			case 5:
				guideDetails.Description.Storage = fmt.Sprint(guideDetails.Description.Storage, text)
			case 6:
				guideDetails.Description.GPU = fmt.Sprint(guideDetails.Description.GPU, text)
			case 7:
				guideDetails.Description.Case = fmt.Sprint(guideDetails.Description.Case, text)
			case 8:
				guideDetails.Description.PSU = fmt.Sprint(&guideDetails.Description.PSU, text)
			}
		}
	}

	commentNumber, _  := scraper.Instance.FindElement(selenium.ByCSSSelector, "#comments")
	comText, _ := commentNumber.Text()

	guideDetails.NumberOfComments, _ = strconv.Atoi(comText)

	_ = json.NewEncoder(w).Encode(guideDetails)

	b, _ := json.Marshal(guideDetails)
	cache.Put(path, b)
}