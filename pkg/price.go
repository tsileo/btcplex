package btcplex

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Return last USD price from BitcoinAverage API
func GetLastBitcoinPrice() (price float64, err error) {
	resp, err := http.Get("https://api.bitcoinaverage.com/ticker/global/USD/last")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	price, err = strconv.ParseFloat(string(body), 10)
	return
}

func main() {
	p, _ := GetLastBitcoinPrice()
	fmt.Printf("%v", p)
}
