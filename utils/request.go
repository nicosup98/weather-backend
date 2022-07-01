package utils

import (
	"log"

	"github.com/valyala/fasthttp"
)

func MakeRequest(method string, url string) fasthttp.Response {
	log.Println(url)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.Add("X-RapidAPI-Key", GetEnviromentVars("RAPID_API_KEY"))
	req.Header.Add("X-RapidAPI-Host", GetEnviromentVars("RAPID_API_HOST"))

	resp := fasthttp.AcquireResponse()

	err := fasthttp.Do(req, resp)

	if err != nil {
		log.Panicln("an error ocurred doing http request: ", err)
	}

	return *resp
}
