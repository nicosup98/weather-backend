package main

import (
	"fmt"
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
		AllowOrigins: "http://localhost",
		AllowHeaders: "Origin, Content-Type, Accept, session_id",
	}))

	server.Get("weather/:city", weather.GetWeather)

	server.Get("autocomplete/:city", weather.GetAutocompletation)

	server.Get("history", weather.GetHistorial)

	server.Listen(fmt.Sprintf(":%s", local_utils.GetEnviromentVars("PORT")))
}
