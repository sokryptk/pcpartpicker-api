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

func GetGuideDetails(w http.ResponseWriter, r *http.Request) {
	path := r.Header.Get("path")
	region := r.Header.Get("region")

	guideDetails, cached := GetDetails(path, region)

	_ = json.NewEncoder(w).Encode(guideDetails)

	if !cached {
		b, _ := json.Marshal(guideDetails)
		cache.Put(path, b)
	}
}

//Returns true if value retrieved from Cache, else false
func GetDetails(path string, region string) (entities.GuideDetails, bool) {
	if data, success := cache.RetrieveCache(path); success {
		var db entities.GuideDetails

		_ = json.Unmarshal(data, &db)

		return db, true
	}

	if _, err := scraper.Instance.ExecuteScript(fmt.Sprintf("window.open('%s');", path), nil); err != nil {
		log.Println(err)
	}

	windows, _ := scraper.Instance.WindowHandles()
	_ = scraper.Instance.SwitchWindow(windows[len(windows)-1])

	handle, _ := scraper.Instance.CurrentWindowHandle()
	defer scraper.Instance.SwitchWindow(windows[0])
	defer scraper.Instance.CloseWindow(handle)

	var guideDetails entities.GuideDetails

	err := scraper.Instance.WaitWithTimeout(func(wd selenium.WebDriver) (b bool, err error) {
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
	pUrl := parse.Parser{Region: region}
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

	commentNumber, _ := scraper.Instance.FindElement(selenium.ByCSSSelector, "#comments")
	comText, _ := commentNumber.Text()

	guideDetails.NumberOfComments, _ = strconv.Atoi(comText)

	return guideDetails, false
}
