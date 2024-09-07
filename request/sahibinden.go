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

var topPageNumber int
var sahibindenLimit int = 0
var SahinindenProducts []settings.Product

func PageSahibinden(parameter string, min int, max int, descFilter string, show bool, output string, ascending bool, descending bool, limit int) {
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
	} else if min == -1 && max == -1 {
		url = "https://www.sahibinden.com/ikinci-el-ve-sifir-alisveris?query_text_mf=" + parameter + "&query_text=" + parameter
	}

	service, err := selenium.NewChromeDriverService("C:\\WebDriver\\chromedriver.exe", 4444)

	settings.ErrorHandler(err)

	defer service.Stop()
	customUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36"
	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{
		"--user-agent=" + customUserAgent,
		"--headless",
	}})

	driver, err := selenium.NewRemote(caps, "")

	settings.ErrorHandler(err)

	err = driver.Get(url)

	settings.ErrorHandler(err)

	naviUl, err := driver.FindElement(selenium.ByCSSSelector, "ul.pageNaviButtons")

	if err != nil {
		topPageNumber = 1
	} else {
		aTag, err := naviUl.FindElements(selenium.ByTagName, "a")

		settings.ErrorHandler(err)

		for _, v := range aTag {
			text, err := v.Text()

			settings.ErrorHandler(err)

			if integerText, err := strconv.Atoi(text); err == nil {
				if integerText > topPageNumber {
					topPageNumber = integerText
				}
			}
		}
	}

	driver.Close()
	Sahibinden(parameter, min, max, descFilter, show, ascending, descending, limit)

	if output != "" {
		SahibindenOutput(output, SahinindenProducts)
	}
}

func Sahibinden(parameter string, min int, max int, descFilter string, show bool, ascending bool, descending bool, limit int) {

	var baseUrl string = "https://www.sahibinden.com/ikinci-el-ve-sifir-alisveris"

	for i := 1; i <= topPageNumber; i++ {
		var url string
		if i != 1 {

			if min == -1 && max == -1 {
				url = baseUrl + "?pagingOffset=" + strconv.Itoa((i*20)-20) + "&query_text_mf=" + parameter + "&query_text=" + parameter
			} else if min != -1 && max == -1 {
				url = baseUrl + "?pagingOffset=" + strconv.Itoa((i*20)-20) + "&price_min=" + strconv.Itoa(min) + "&query_text_mf=" + parameter + "&query_text=" + parameter
			} else if min == -1 && max != -1 {
				url = baseUrl + "?pagingOffset=" + strconv.Itoa((i*20)-20) + "&query_text_mf=" + parameter + "&price_max=" + strconv.Itoa(max) + "&query_text=" + parameter
			} else if min != -1 && max != -1 {
				url = baseUrl + "?pagingOffset=" + strconv.Itoa((i*20)-20) + "&price_min=" + strconv.Itoa(min) + "&query_text_mf=" + parameter + "&price_max=" + strconv.Itoa(max) + "&query_text=" + parameter
			}
		} else {
			url = baseUrl + "?query_text_mf=" + parameter + "&query_text=" + parameter
			if min == -1 && max == -1 {
				url = baseUrl + "?" + "query_text_mf=" + parameter + "&query_text=" + parameter
			} else if min != -1 && max == -1 {
				url = baseUrl + "?" + "price_min=" + strconv.Itoa(min) + "&query_text_mf=" + parameter + "&query_text=" + parameter
			} else if min == -1 && max != -1 {
				url = baseUrl + "?" + "query_text_mf=" + parameter + "&price_max=" + strconv.Itoa(max) + "&query_text=" + parameter
			} else if min != -1 && max != -1 {
				url = baseUrl + "?" + "price_min=" + strconv.Itoa(min) + "&query_text_mf=" + parameter + "&price_max=" + strconv.Itoa(max) + "&query_text=" + parameter
			}
		}

		if ascending {
			url = url + "&sorting=price_asc"
		} else if descending {
			url = url + "&sorting=price_desc"
		}

		service, err := selenium.NewChromeDriverService("C:\\WebDriver\\chromedriver.exe", 4444)

		settings.ErrorHandler(err)

		defer service.Stop()
		customUserAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/115.0.5790.170 Safari/537.36"
		caps := selenium.Capabilities{}
		caps.AddChrome(chrome.Capabilities{Args: []string{
			"--headless",
			"--user-agent=" + customUserAgent,
		}})

		driver, err := selenium.NewRemote(caps, "")

		settings.ErrorHandler(err)

		fmt.Println("İstek atılıyor ---> ", color.MagentaString(url))
		err = driver.Get(url)

		settings.ErrorHandler(err)

		resultItem, err := driver.FindElements(selenium.ByCSSSelector, "tr.searchResultsItem")

		settings.ErrorHandler(err)

		for _, v := range resultItem {

			text, _ := v.Text()

			if text != "" {
				thumbnail, err := v.FindElement(selenium.ByCSSSelector, "td.searchResultsLargeThumbnail")

				settings.ErrorHandler(err)

				tdPrice, err := v.FindElement(selenium.ByCSSSelector, "td.searchResultsPriceValue")

				settings.ErrorHandler(err)

				priceTag, err := tdPrice.FindElement(selenium.ByTagName, "span")
				settings.ErrorHandler(err)

				price, err := priceTag.Text()

				settings.ErrorHandler(err)
				aTag, err := thumbnail.FindElement(selenium.ByTagName, "a")

				settings.ErrorHandler(err)

				href, _ := aTag.GetAttribute("href")
				title, _ := aTag.GetAttribute("title")

				splitString := strings.Split(price, " ")

				money := splitString[0]
				var newMoney int

				if strings.Contains(money, ".") {
					newMoney, err = strconv.Atoi(strings.ReplaceAll(money, ".", ""))

					settings.ErrorHandler(err)
				} else {
					newMoney, err = strconv.Atoi(money)

					settings.ErrorHandler(err)
				}

				if descFilter == "" && min == 0 && max == 0 {
					newProduct := settings.Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						settings.Show(newProduct)
					}
				} else if descFilter != "" && strings.Contains(title, descFilter) && min == 0 && max == 0 {
					newProduct := settings.Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						settings.Show(newProduct)
					}
				} else if descFilter == "" && ((min != 0 && min <= newMoney) || (max != 0 && newMoney <= max)) {
					newProduct := settings.Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						settings.Show(newProduct)
					}
				} else if descFilter != "" && strings.Contains(title, descFilter) && ((min != 0 && min <= newMoney) || (max != 0 && newMoney <= max)) {
					newProduct := settings.Product{href, title, price}
					SahinindenProducts = append(SahinindenProducts, newProduct)
					if !show {
						settings.Show(newProduct)
					}
				}
			}
		}
		driver.Close()
	}
}

func SahibindenOutput(output string, sahibinden []settings.Product) {

	file, err := os.Create(output)

	settings.ErrorHandler(err)

	for _, v := range sahibinden {
		file.WriteString("---Sahibinden---")
		file.WriteString(v.URL)
		file.WriteString(v.Desc)
		file.WriteString(v.Price)
		file.WriteString("-------------------------------------------------")
	}
}
