package redis

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/go-redis/redis"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var redisDb redisClient

func TestMain(m *testing.M) {

	monkey.Patch(newClient, func() redisClient {
		client := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", "192.168.1.167", 6379),
			Password: "",
			DB:       0,
		})
		return redisClient{client: client}
	})
	m.Run()
}

func TestRedis(t *testing.T) {
	redisDb = newClient()
	Convey("testing redis set", t, func() {
		So(redisDb.Set("test", "1111", 0), ShouldBeNil)
		Convey("testing redis get", func() {
			val, err := redisDb.Get("test")
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "1111")
		})
	})

}
