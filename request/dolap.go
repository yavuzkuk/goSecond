package request

import (
	"fmt"
	"go2Second/settings"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var topPage int = 0
var dolapLimit int = 0
var productMap map[string]string
var dolapProducts []settings.Product

// var urls map[string]string

func PageNumberDolap(parameter string, min int, max int, output string, descFilter string, show bool, ascending bool, descending bool, limit int) {
	var order string
	url := "https://dolap.com/ara?q=" + parameter

	if ascending {
		order = "&sira=artan-fiyat"
		url = url + order
	} else if descending {
		order = "&sira=azalan-fiyat"
		url = url + order
	}

	fmt.Println("İstek atılıyor ----> ", url)

	service, err := selenium.NewChromeDriverService("C:\\WebDriver\\chromedriver.exe", 4444)

	settings.ErrorHandler(err)

	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		// "--headless",
	}})

	driver, err := selenium.NewRemote(caps, "")

	settings.ErrorHandler(err)

	err = driver.Get(url)
	settings.ErrorHandler(err)

	ul, err := driver.FindElements(selenium.ByCSSSelector, "ul.pagination.other")

	settings.ErrorHandler(err)

	var paginationUl selenium.WebElement
	if len(ul) > 1 {
		paginationUl = ul[1]
	} else {
		paginationUl = ul[0]
	}

	liElements, err := paginationUl.FindElements(selenium.ByTagName, "a")

	settings.ErrorHandler(err)

	for _, v := range liElements {
		text, err := v.Text()
		settings.ErrorHandler(err)

		if text != "" {
			textInteger, err := strconv.Atoi(text)

			settings.ErrorHandler(err)

			if textInteger > topPage {
				topPage = textInteger
			}
		}
	}

	GetProduct(parameter, min, max, output, descFilter, show, driver, topPage, ascending, descending)

	if output != "" {
		Output(output, dolapProducts)
	}
}

func GetProduct(parameter string, min int, max int, output string, descFilter string, show bool, driver selenium.WebDriver, intPage int, ascending bool, descending bool) {

	url := "https://dolap.com/ara?q=" + parameter

	if ascending {
		url = url + "&sira=artan-fiyat"
	} else if descending {
		url = url + "&sira=azalan-fiyat"
	}

	for i := 1; i < intPage; i++ {
		productMap = make(map[string]string)
		newUrl := url + "&sayfa=" + strconv.Itoa(i)
		fmt.Println("İstek atılıyor --->", color.GreenString(newUrl))

		fmt.Println("URL ------> ", newUrl)

		err := driver.Get(newUrl)

		settings.ErrorHandler(err)

		colHolders, err := driver.FindElements(selenium.ByCSSSelector, "div.col-holder")
		settings.ErrorHandler(err)

		for _, holder := range colHolders {
			imgBlock, err := holder.FindElement(selenium.ByCSSSelector, "div.img-block")
			settings.ErrorHandler(err)

			aTag, err := imgBlock.FindElement(selenium.ByTagName, "a")
			settings.ErrorHandler(err)

			href, err := aTag.GetAttribute("href")
			settings.ErrorHandler(err)

			priceSpan, err := holder.FindElement(selenium.ByCSSSelector, "span.price")
			settings.ErrorHandler(err)

			priceText, err := priceSpan.Text()
			settings.ErrorHandler(err)
			productMap[href] = priceText
		}
		ProductDetail(productMap, min, max, descFilter, driver)
	}

}

func ProductDetail(product map[string]string, min int, max int, descFilter string, driver selenium.WebDriver) {

	for url, price := range product {

		err := driver.Get(url)

		settings.ErrorHandler(err)
		elementDesc, err := driver.FindElement(selenium.ByClassName, "remarks-block")

		settings.ErrorHandler(err)
		desc, err := elementDesc.FindElement(selenium.ByTagName, "p")

		settings.ErrorHandler(err)

		descText, err := desc.Text()

		settings.ErrorHandler(err)

		splitInteger := strings.Split(price, " ")
		money := splitInteger[0]

		if strings.Contains(money, ".") {
			money = strings.ReplaceAll(money, ".", "")
		}

		integerPrice, err := strconv.Atoi(money)
		settings.ErrorHandler(err)

		if descFilter == "" && min == -1 && max == -1 {
			newProduct := settings.Product{url, descText, price}
			dolapProducts = append(dolapProducts, newProduct)
			settings.Show(newProduct)
		} else if descFilter != "" && strings.Contains(descText, descFilter) && min == -1 && max == -1 {
			newProduct := settings.Product{url, descText, price}
			dolapProducts = append(dolapProducts, newProduct)
			settings.Show(newProduct)
		} else if descFilter == "" && (min != -1 && min <= integerPrice) || (max != -1 && integerPrice <= max) {
			newProduct := settings.Product{url, descText, price}
			dolapProducts = append(dolapProducts, newProduct)
			settings.Show(newProduct)
		} else if descFilter != "" && strings.Contains(descText, descFilter) && ((min != -1 && min <= integerPrice) || (max != -1 && integerPrice <= max)) {
			newProduct := settings.Product{url, descText, price}
			dolapProducts = append(dolapProducts, newProduct)
			settings.Show(newProduct)
		}
	}
}

func Output(output string, prodcuts []settings.Product) {

	file, err := os.Create(output)

	settings.ErrorHandler(err)

	for _, v := range prodcuts {
		file.WriteString(v.URL)
		file.WriteString(v.Desc)
		file.WriteString(v.Price)
		file.WriteString("-------------------------------------------------")
	}

}
