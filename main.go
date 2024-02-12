package main

import (
	"fmt"
	"log"
	redis_session "weather-app/api/session"
	local_utils "weather-app/api/utils"
	"weather-app/api/weather"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	server := fiber.New()
	redis_session.New()
	server.Use(cors.New(cors.Config{
		AllowOrigins:  "https://weather-app-eight-steel.vercel.app/",
		AllowHeaders:  "Origin, Content-Type, Accept, session_id",
		ExposeHeaders: "session_id",
	}))

	server.Get("weather/:city", weather.GetWeather)

	server.Get("autocomplete/:city", weather.GetAutocompletation)

	server.Get("history", weather.GetHistorial)
	server.Get("getToken", redis_session.GetToken)
	server.Delete("deleteToken", redis_session.DeleteSessionToken)

	log.Fatal(server.Listen(fmt.Sprintf(":%s", local_utils.GetEnviromentVars("PORT"))))
}
