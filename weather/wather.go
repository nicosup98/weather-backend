package weather

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"

	redis_session "weather-app/api/session"
	local_utils "weather-app/api/utils"
)

var (
	store = redis_session.Store
)

func GetWeather(c *fiber.Ctx) error {
	forecast := c.Query("daysForecast", "")
	city := c.Params("city")

	var rawUrl string
	var typeSearch string
	params := map[string]string{
		"q": city,
	}

	if forecast == "" {
		rawUrl = "https://weatherapi-com.p.rapidapi.com/current.json"
		typeSearch = "current"

	} else {
		rawUrl = "https://weatherapi-com.p.rapidapi.com/forecast.json"
		params["days"] = forecast
		typeSearch = "forecast"
	}

	stdUrl, err := local_utils.UrlParser(rawUrl, params)

	if err != nil {
		log.Panicln("an error ocurred parsin url: ", err)
	}

	resp := local_utils.MakeRequest(http.MethodGet, stdUrl)
	defer fasthttp.ReleaseResponse(&resp)

	body := string(resp.Body())
	code := resp.StatusCode()

	c.Response().SetStatusCode(code)

	sess, err := store.Get(c)

	if err != nil {
		log.Panicln("an error ocurred getting store: ", err)
	}

	timeStr := time.Now().String()

	historyData := map[string]interface{}{
		"type": typeSearch,
		"data": body,
	}

	sess.Set(timeStr, historyData)

	if err := sess.Save(); err != nil {
		log.Panicln("an error ocurred saving the session: ", err)
	}

	return c.JSON(body)

}

func GetAutocompletation(c *fiber.Ctx) error {
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
}

func GetHistorial(c *fiber.Ctx) error {
	lastSearchs := c.Query("lastestSearches", "3")

	sess, err := store.Get(c)

	if err != nil {
		log.Panicln("an error ocurred calling session store: ", err)
	}

	keys := sess.Keys()

	sort.Slice(keys, func(i, j int) bool {
		a, err := time.Parse("15:04:05", keys[i])
		b, err := time.Parse("15:04:05", keys[j])

		if err != nil {
			log.Panicln("an error ocurred parsing time: ", err)
		}

		return b.Before(a)

	})
	lastSearchesN, err := strconv.Atoi(lastSearchs)

	if err != nil {
		log.Panicln("an error ocurred parsing lastSearches: ", err)
	}
	results := []interface{}{}
	for _, v := range keys[:lastSearchesN] {
		r := sess.Get(v)

		results = append(results, r)
	}
	c.Response().SetStatusCode(200)

	return c.JSON(results)

}
