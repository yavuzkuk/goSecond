package settings

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Product struct {
	URL   string
	Desc  string
	Price string
}

func ErrorHandler(err error) {
	if err != nil {
		fmt.Println("Error ----> ", err)
	}
}

func Show(product Product) {

	var unicodeEmoji string
	if strings.Contains(product.URL, "sahibinden") {
		unicodeEmoji = "\U0001F7E8"
		fmt.Println(unicodeEmoji, color.YellowString("Sahibinden"))
		fmt.Println(color.YellowString(product.Desc))

	} else {
		unicodeEmoji = "\U0001F7E9"
		fmt.Println(unicodeEmoji, color.GreenString("Dolap"))
		fmt.Println(color.GreenString(product.Desc))

	}

	fmt.Println(product.URL)
	fmt.Println(product.Price)
	fmt.Println("---------------------------------------")

}
