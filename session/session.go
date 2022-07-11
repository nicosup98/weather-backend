package sesssion

import (
	"encoding/gob"
	"log"
	local_utils "weather-app/api/utils"

	"github.com/gofiber/fiber/v2"
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
		Storage:   storage,
		KeyLookup: "header:session_id",
	})

	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}

func GetToken(c *fiber.Ctx) error {

	sess, err := Store.Get(c)

	if err != nil {
		log.Panicln("an ocurred getting session: ", err)

	}

	return c.Status(200).SendString(sess.ID())

}

func DeleteSessionToken(c *fiber.Ctx) error {
	sess, err := Store.Get(c)

	if err != nil {
		log.Panicln("an errro ocurred getting session: ", err)
	}

	err = sess.Destroy()

	if err != nil {
		log.Panicln("an error ocurred deleteing session: ", err)
	}

	return c.Status(200).SendString("session deleted")

}
