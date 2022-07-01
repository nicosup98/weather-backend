package main

import (
	"fmt"
	"log"
	"net/http"
	local_utils "weater-app/api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	server := fiber.New()

	server.Get("weather/:city", func(c *fiber.Ctx) error {
		forecast := c.Query("daysForecast", "")
		city := c.Params("city")

		var rawUrl string
		params := map[string]string{
			"q": city,
		}

		if forecast == "" {
			rawUrl = "https://weatherapi-com.p.rapidapi.com/current.json"

		} else {
			rawUrl = "https://weatherapi-com.p.rapidapi.com/forecast.json"
			params["days"] = forecast
		}

		stdUrl, err := local_utils.UrlParser(rawUrl, params)

		if err != nil {
			log.Panicln("an error ocurred parsin url: ", err)
		}

		resp := local_utils.MakeRequest(http.MethodGet, stdUrl)
		defer fasthttp.ReleaseResponse(&resp)

		body := resp.Body()
		code := resp.StatusCode()

		c.Response().SetStatusCode(code)

		return c.JSON(string(body))

	})

	server.Get("autocomplete/:city", func(c *fiber.Ctx) error {
		URL, err := local_utils.UrlParser("https://weatherapi-com.p.rapidapi.com/search.json", map[string]string{"q": c.Params("city")})

		if err != nil {
			log.Panicln("an error ocurred parsing url: ", err)
		}

		resp := local_utils.MakeRequest(http.MethodGet, URL)
		defer fasthttp.ReleaseResponse(&resp)

		body := resp.Body()
		code := resp.StatusCode()

		c.Response().SetStatusCode(code)
		return c.JSON(string(body))
	})

	server.Listen(fmt.Sprintf(":%s", local_utils.GetEnviromentVars("PORT")))
}
