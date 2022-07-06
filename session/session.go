package sesssion

import (
	"encoding/gob"
	local_utils "weather-app/api/utils"

	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis"
)

var (
	Store *session.Store
)

func New() {
	storage := redis.New(redis.Config{
		URL:   local_utils.GetEnviromentVars("REDIS_URL"),
		Reset: false,
	})
	Store = session.New(session.Config{
		Storage:      storage,
		KeyLookup:    "header:session_id",
		CookieDomain: "http://localhost:300",
	})

	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}
