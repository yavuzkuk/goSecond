package request

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var topPageNumber int
var pageUrls = make(map[string]string)
var SahinindenProducts []Product

func PageSahibinden(parameter string, min int, max int, descFilter string, show bool, output string) {
	var url string = "https://www.sahibinden.com/ikinci-el-ve-sifir-alisveris?query_text_mf=" + parameter + "&query_text=" + parameter

	if min != -1 && max != -1 {
		minInteger := strconv.Itoa(min)
		maxInteger := strconv.Itoa(max)
		url = "https://www.sahibinden.com/ikinci-el-ve-sifir-alisveris?price_min=" + minInteger + "&query_text_mf=" + parameter + "&price_max=" + maxInteger + "&query_text=" + parameter
	} else if min != -1 && max == -1 {
		minInteger := strconv.Itoa(min)
		url = "https://www.sahibinden.com/ikinci-el-ve-sifir-alisveris?price_min=" + minInteger + "&query_text_mf=" + parameter + "&query_text=" + parameter
	} else if min == -1 && max != -1 {
		maxInteger := strconv.Itoa(max)
		url = "https://www.sahibinden.com/ikinci-el-ve-sifir-alisveris?" + "query_text_mf=" + parameter + "&price_max=" + maxInteger + "&query_text=" + parameter
	}

	pageUrls["1"] = url
	service, err := selenium.NewChromeDriverService("C:\\WebDriver\\chromedriver.exe", 4444)

	if err != nil {
		fmt.Println("Service error --> ", err)
	}

	defer service.Stop()
	// proxyServerURL := "213.157.6.50"
	customUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36"
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--user-agent=" + customUserAgent,
		"--headless",
	}})

	driver, err := selenium.NewRemote(caps, "")

	if err != nil {
		fmt.Println("Newremote error --> ", err)
	}

	err = driver.Get(url)

	if err != nil {
		fmt.Println("Error --> ", err)
	}

	naviUl, err := driver.FindElement(selenium.ByCSSSelector, "ul.pageNaviButtons")

	if err != nil {
		topPageNumber = 1
	} else {
		aTag, err := naviUl.FindElements(selenium.ByTagName, "a")

		if err != nil {
			fmt.Println("A tag error --> ", err)
		}

		for _, v := range aTag {
			text, err := v.Text()

			if err != nil {
				fmt.Println("Text error --> ", err)
			}

			if integerText, err := strconv.Atoi(text); err == nil {
				href, _ := v.GetAttribute("href")
				pageUrls[text] = href
				if integerText > topPageNumber {
					topPageNumber = integerText
				}
			}
		}
	}

	driver.Close()
	Sahibinden(parameter, min, max, descFilter, show)

	if output != "" {
		SahibindenOutput(output, SahinindenProducts)
	}
}

func Sahibinden(parameter string, min int, max int, descFilter string, show bool) {
	for _, url := range pageUrls {
		service, err := selenium.NewChromeDriverService("C:\\WebDriver\\chromedriver.exe", 4444)

		if err != nil {
			fmt.Println("Service error --> ", err)
		}

		defer service.Stop()
		customUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36"
		caps := selenium.Capabilities{}
		caps.AddChrome(chrome.Capabilities{Args: []string{
			"--headless",
			"--user-agent=" + customUserAgent,
		}})

		driver, err := selenium.NewRemote(caps, "")

		if err != nil {
			fmt.Println("Newremote error --> ", err)
		}

		fmt.Println("İstek atılıyor ---> ", color.MagentaString(url))
		err = driver.Get(url)

		if err != nil {
			fmt.Println("Get error --->", err)
		}

		file, _ := os.Create("index.html")

		html, _ := driver.PageSource()
		file.WriteString(html)

		resultItem, err := driver.FindElements(selenium.ByCSSSelector, "tr.searchResultsItem")

		if err != nil {
			fmt.Println("Result item error --> ", err)
		}

		for _, v := range resultItem {

			text, _ := v.Text()

			if text != "" {
				thumbnail, err := v.FindElement(selenium.ByCSSSelector, "td.searchResultsLargeThumbnail")

				if err != nil {
					fmt.Println("Thumbnail error ---> ", err)
				}
				tdPrice, err := v.FindElement(selenium.ByCSSSelector, "td.searchResultsPriceValue")

				if err != nil {
					fmt.Println("Thumbnail error ---> ", err)
				}

				priceTag, err := tdPrice.FindElement(selenium.ByTagName, "span")
				if err != nil {
					fmt.Println("Span error --> ", err)
				}

				price, err := priceTag.Text()

				if err != nil {
					fmt.Println("Span error --> ", err)
				}
				aTag, err := thumbnail.FindElement(selenium.ByTagName, "a")

				if err != nil {
					fmt.Println("A tag error --> ", err)
				}

				href, _ := aTag.GetAttribute("href")
				title, _ := aTag.GetAttribute("title")

				splitString := strings.Split(price, " ")

				money := splitString[0]
				moneyType := splitString[1]
				var newMoney int

				if strings.Contains(money, ".") {
					newMoney, err = strconv.Atoi(strings.ReplaceAll(money, ".", ""))

					if err != nil {
						fmt.Println("Money error --> ", err)
					}
				} else {
					newMoney, err = strconv.Atoi(money)

					if err != nil {
						fmt.Println("Money error --> ", err)
					}
				}

				price = strconv.Itoa(newMoney)
				if descFilter == "" && min == 0 && max == 0 {
					newProduct := Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						SahibindenShow(newProduct, moneyType)
					}
				} else if descFilter != "" && strings.Contains(title, descFilter) && min == 0 && max == 0 {
					newProduct := Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						SahibindenShow(newProduct, moneyType)
					}
				} else if descFilter == "" && ((min != 0 && min <= newMoney) || (max != 0 && newMoney <= max)) {
					newProduct := Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						SahibindenShow(newProduct, moneyType)
					}
				} else if descFilter != "" && strings.Contains(title, descFilter) && ((min != 0 && min <= newMoney) || (max != 0 && newMoney <= max)) {
					newProduct := Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						SahibindenShow(newProduct, moneyType)
					}
				}
			}
		}
		driver.Close()
	}
}

func SahibindenShow(product Product, moneyType string) {
	fmt.Println(color.YellowString(product.Desc))
	fmt.Println(product.URL)
	fmt.Println(product.Price, moneyType)
	fmt.Println("---------------------------------------")
}

func SahibindenOutput(output string, sahibinden []Product) {

	file, err := os.Create(output)

	if err != nil {
		fmt.Println("File create error --> ", err)
	}

	for _, v := range sahibinden {
		file.WriteString("---Sahibinden---")
		file.WriteString(v.URL)
		file.WriteString(v.Desc)
		file.WriteString(v.Price)
		file.WriteString("-------------------------------------------------")
	}

}
