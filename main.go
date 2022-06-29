package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"weater-app/api/model"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

func main() {
	server := fiber.New()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("some error ocurred parsing .env %s", err)
	}

	server.Get("current_weather/:city", func(c *fiber.Ctx) error {
		params := model.WeatherRequest{
			Ip:      c.Query("ip"),
			Region:  c.Query("region"),
			Country: c.Query("country"),
		}
		city := c.Params("city")
		rawUrl := "https://weatherapi-com.p.rapidapi.com/current.json"
		stdUrl, err := url.Parse(rawUrl)

		if err != nil {
			log.Panicln("an error ocurred parsin url: ", err)
		}

		query := stdUrl.Query()

		query.Add("q", fmt.Sprintf("%s %s", city, params.GetPreciseLocation()))

		stdUrl.RawQuery = query.Encode()

		resp := makeRequest(http.MethodGet, stdUrl.String())
		defer fasthttp.ReleaseResponse(&resp)

		body := resp.Body()
		// copy(body, resp.Body())
		code := resp.StatusCode()

		log.Println("body", string(body))

		c.Response().SetStatusCode(code)

		return c.JSON(string(body))

	})

	server.Listen(":5000")
}

func makeRequest(method string, url string) fasthttp.Response {
	log.Println(url)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.Add("X-RapidAPI-Key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", os.Getenv("RAPID_API_HOST"))

	resp := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, resp)

	if err != nil {
		log.Panicln("an error ocurred doing http request: ", err)
	}

	return *resp
}
