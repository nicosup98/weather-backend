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

	server.Get("current_weather/:city", func(c *fiber.Ctx) error {
		city := c.Params("city")
		rawUrl := "https://weatherapi-com.p.rapidapi.com/current.json"
		stdUrl, err := local_utils.UrlParser(rawUrl, city)

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
		URL, err := local_utils.UrlParser("https://weatherapi-com.p.rapidapi.com/search.json", c.Params("city"))

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
