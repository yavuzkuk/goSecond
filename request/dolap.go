package request

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var topPage int = 0
var productMap map[string]string
var Products []Product

type Product struct {
	URL   string
	Desc  string
	Price string
}

func PageNumberDolap(parameter string, min int, max int, output string, descFilter string, show bool) {
	url := "https://dolap.com/ara?q=" + parameter

	service, err := selenium.NewChromeDriverService("C:\\WebDriver\\chromedriver.exe", 4444)

	if err != nil {
		log.Fatal("Error:", err)
	}

	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--headless",
	}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error:", err)
	}

	err = driver.Get(url)
	if err != nil {
		log.Fatal("Error:", err)
	}

	ul, err := driver.FindElements(selenium.ByCSSSelector, "ul.pagination.other")

	if err != nil {
		fmt.Println("Error ---> ", err)
	}
	var paginationUl selenium.WebElement
	if len(ul) > 1 {
		paginationUl = ul[1]
	} else {
		paginationUl = ul[0]
	}

	liElements, err := paginationUl.FindElements(selenium.ByTagName, "a")

	if err != nil {
		fmt.Println("Li elements error --> ", err)
	}

	for _, v := range liElements {
		text, err := v.Text()
		if err != nil {
			fmt.Println("Error --> ", err)
		}

		if text != "" {
			textInteger, err := strconv.Atoi(text)

			if err != nil {
				fmt.Println("Error --> ", err)
			}

			if textInteger > topPage {
				topPage = textInteger
			}
		}
	}

	GetProduct(parameter, min, max, output, descFilter, show, driver, topPage)

	if output != "" {
		Output(output, Products)
	}
}

func GetProduct(parameter string, min int, max int, output string, descFilter string, show bool, driver selenium.WebDriver, intPage int) {

	productMap = make(map[string]string)

	fmt.Println(intPage)
	for i := 1; i < intPage; i++ {
		url := "https://dolap.com/ara?q=" + parameter + "&sayfa=" + strconv.Itoa(i)

		fmt.Println("İstek atılıyor --->", color.GreenString(url))

		err := driver.Get(url)

		if err != nil {
			fmt.Println("Get error --> ", err)
		}

		colHolders, err := driver.FindElements(selenium.ByCSSSelector, "div.col-holder")
		if err != nil {
			log.Fatalf("Hedef etiket bulunamadı: %v", err)
		}

		for _, holder := range colHolders {
			imgBlock, err := holder.FindElement(selenium.ByCSSSelector, "div.img-block")
			if err != nil {
				fmt.Println("img-block bulunamadı:", err)
				continue
			}

			aTag, err := imgBlock.FindElement(selenium.ByTagName, "a")
			if err != nil {
				fmt.Println("a etiketi bulunamadı:", err)
				continue
			}

			href, err := aTag.GetAttribute("href")
			if err != nil {
				fmt.Println("Href bulunamadı:", err)
			}

			priceSpan, err := holder.FindElement(selenium.ByCSSSelector, "span.price")
			if err != nil {
				fmt.Println("Fiyat bulunamadı:", err)
				continue
			}

			priceText, err := priceSpan.Text()
			if err != nil {
				fmt.Println("Fiyat metni alınamadı:", err)
			}
			productMap[href] = priceText
		}
	}
	ProductDetail(productMap, min, max, descFilter, driver)

}

func ProductDetail(product map[string]string, min int, max int, descFilter string, driver selenium.WebDriver) {

	for url, price := range product {
		err := driver.Get(url)

		if err != nil {
			fmt.Println("Get error --> ", err)
		}

		elementDesc, err := driver.FindElement(selenium.ByClassName, "remarks-block")

		if err != nil {
			fmt.Println("Desc error --> ", err)
		}

		desc, err := elementDesc.FindElement(selenium.ByTagName, "p")

		if err != nil {
			fmt.Println("Element error --> ", err)
		}

		descText, err := desc.Text()

		if err != nil {
			fmt.Println("Desc error --> ", err)
		}

		integerPrice, _ := strconv.Atoi(price)

		if descFilter == "" && min == -1 && max == -1 {
			newProduct := Product{url, descText, price}
			Products = append(Products, newProduct)
			Show(newProduct)
		} else if descFilter != "" && strings.Contains(descText, descFilter) && min == -1 && max == -1 {
			newProduct := Product{url, descText, price}
			Products = append(Products, newProduct)
			Show(newProduct)
		} else if descFilter == "" && ((min != -1 && min <= integerPrice) || (max != -1 && integerPrice <= max)) {
			newProduct := Product{url, descText, price}
			Products = append(Products, newProduct)
			Show(newProduct)
		} else if descFilter != "" && strings.Contains(descText, descFilter) && ((min != -1 && min <= integerPrice) || (max != -1 && integerPrice <= max)) {
			newProduct := Product{url, descText, price}
			Products = append(Products, newProduct)
			Show(newProduct)
		}

	}
}

func Show(products Product) {

	fmt.Println(color.GreenString(products.URL))
	fmt.Println(products.Desc)
	fmt.Println(products.Price)
	fmt.Println("---------------------------------------")
}

func Output(output string, prodcuts []Product) {

	file, err := os.Create(output)

	if err != nil {
		fmt.Println("File create error --> ", err)
	}

	for _, v := range prodcuts {
		file.WriteString(v.URL)
		file.WriteString(v.Desc)
		file.WriteString(v.Price)
		file.WriteString("-------------------------------------------------")
	}

}
