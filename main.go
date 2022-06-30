package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"weater-app/api/model"
	local_utils "weater-app/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	server := fiber.New()

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
		code := resp.StatusCode()

		c.Response().SetStatusCode(code)

		return c.JSON(string(body))

	})

	server.Listen(fmt.Sprintf(":%s", local_utils.GetEnviromentVars("PORT")))
}

func makeRequest(method string, url string) fasthttp.Response {
	log.Println(url)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.Add("X-RapidAPI-Key", local_utils.GetEnviromentVars("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", local_utils.GetEnviromentVars("RAPID_API_HOST"))

	resp := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, resp)

	if err != nil {
		log.Panicln("an error ocurred doing http request: ", err)
	}

	return *resp
}
