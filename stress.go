package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"

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
	if err != nil {
		fmt.Println("Error -> ", err)
	}

	wg := sync.WaitGroup{}

	for i := 0; i < intLimit; i++ {
		wg.Add(1)
		go callApi(i, &wg)
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
	wg.Wait()
}

func callApi(i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("test", i)
	resp, err := http.PostForm("https://api.pateh.com/api/safar-market-search",
		url.Values{
			"from":          {"THR"},
			"to":            {"KIH"},
			"departureDate": {"2023-06-08"},
			"adult":         {"1"},
			"child":         {"0"},
			"infant":        {"0"},
		})
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	_ = body
	//fmt.Println(string(body))
}
