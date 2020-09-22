package scraper

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var Instance selenium.WebDriver
var Service *selenium.Service

const (
	chromeDriverPath = "ext/chromedriver"
	port             = 6969
)

func init() {
	var opts []selenium.ServiceOption
	var err error

	Service, err = selenium.NewChromeDriverService(chromeDriverPath, port, opts...)
	if err != nil {
		panic(err) //Haha yes
	}

	caps := selenium.Capabilities{"browserName": "chrome"}
	chromeCaps := chrome.Capabilities{
		Args: []string{
			"--headless",
		},
		Path: "/app/.apt/usr/bin/google-chrome",
	}

	caps.AddChrome(chromeCaps)
	Instance, err = selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
}
