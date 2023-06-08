package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/safarmarket", safarmarketTest)
	return r
}

func safarmarketTest(c *gin.Context) {
	limit := c.PostForm("limit")
	intLimit, err := strconv.Atoi(limit)
	days := c.PostForm("days")
	intDays, err := strconv.Atoi(days)

	if err != nil {
		fmt.Println("Error -> ", err)
	}

	wg := sync.WaitGroup{}

	domesticAirports := [3]string{"THR", "MHD", "KIH"}
	internationalAirports := [3]string{"IKA", "MHD", "ISF"}

	currentTime := time.Now()

	for z := 0; z < intDays; z++ {
		date := currentTime.AddDate(0, 0, z)
		for i := 0; i < intLimit; {
			for n := 0; n < len(domesticAirports); n++ {
				for nn := 0; nn < len(domesticAirports); nn++ {
					i++
					wg.Add(1)
					fmt.Println("tt -> ", n, domesticAirports[n], "->", domesticAirports[nn])
					go callApi(i, &wg, "https://api.pateh.com", domesticAirports[n], domesticAirports[nn], date.Format("2006-1-2"))
				}
			}
			for m := 0; m < len(internationalAirports); m++ {
				for mm := 0; mm < len(domesticAirports); mm++ {
					i++
					wg.Add(1)
					fmt.Println("tt -> ", m, internationalAirports[m], "->", internationalAirports[mm])
					go callApi(i, &wg, "https://api.pateh.com", internationalAirports[m], internationalAirports[mm], date.Format("2006-1-2"))
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
	wg.Done()
}

func callApi(i int, wg *sync.WaitGroup, baseUrl string, from string, to string, date string) {
	defer wg.Done()
	//fmt.Println("test", i)
	resp, err := http.PostForm(baseUrl+"/api/safar-market-search",
		url.Values{
			"from":          []string{from},
			"to":            []string{to},
			"departureDate": []string{date},
			"adult":         {"1"},
			"child":         {"0"},
			"infant":        {"0"},
		})
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		//os.Exit(1)
	}
	//defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	_ = resp
	_ = body
	//fmt.Println(string(body))
}
