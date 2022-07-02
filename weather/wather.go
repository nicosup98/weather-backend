package weather

import (
	"encoding/json"
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

	body := resp.Body()
	code := resp.StatusCode()

	c.Response().SetStatusCode(code)

	sess, err := redis_session.Store.Get(c)

	if err != nil {
		log.Panicln("an error ocurred getting store: ", err)
	}

	timeUnix := time.Now().Unix()

	var bodyParsed interface{}

	json.Unmarshal(body, bodyParsed)

	historyData := map[string]interface{}{
		"type": typeSearch,
		"data": bodyParsed,
	}

	historyDataRaw, err := json.Marshal(historyData)

	if err != nil {
		log.Panicln("an error ocurred parsing to json: ", err)
	}

	sess.Set(strconv.FormatInt(timeUnix, 10), string(historyDataRaw))

	if err := sess.Save(); err != nil {
		log.Panicln("an error ocurred saving the session: ", err)
	}

	return c.JSON(string(body))

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

	sess, err := redis_session.Store.Get(c)

	if err != nil {
		log.Panicln("an error ocurred calling session store: ", err)
	}

	keys := sess.Keys()

	sort.Slice(keys, func(i, j int) bool {
		timestampStrA, err := strconv.Atoi(keys[i])
		timestampStrB, err := strconv.Atoi(keys[j])

		if err != nil {
			log.Panicln("an error ocurred parsing unix timestamp: ", err)
		}

		a := time.UnixMilli(int64(timestampStrA))
		b := time.UnixMilli(int64(timestampStrB))

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
	var keysSliced []string

	if len(keys) < lastSearchesN {
		keysSliced = keys
	} else {
		keysSliced = keys[:lastSearchesN]
	}
	for _, v := range keysSliced {
		r := sess.Get(v)

		results = append(results, r)
	}
	c.Response().SetStatusCode(200)

	return c.JSON(results)

}
