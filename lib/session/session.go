package session

import (
	"fmt"
	"gin_websocket/lib/redis"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type session struct {
	sid string
}

const SidLength = 32
const sessionName = "gid"
const LifeTime = 864000 * time.Second

var sidReg = regexp.MustCompile(fmt.Sprintf("[a-z0-9]{%d}", SidLength))

func NewSession(cRequest *http.Request, cResponse gin.ResponseWriter) *session {
	var sid string
	cookie, _ := cRequest.Cookie(sessionName)
	if cookie == nil || !sidReg.MatchString(cookie.Value) {
		sid = genSid()
	} else {
		sid = cookie.Value
	}
	http.SetCookie(cResponse, &http.Cookie{Name: sessionName, Value: sid, Path: "/", HttpOnly: true, Secure: true, Expires: time.Now().Add(LifeTime)})
	return &session{
		sid: sid,
	}
}

func (session *session) GetString(key string) (string, error) {
	return redis.RedisDb.HGet(session.sid, key)
}

func (session *session) Set(key string, value string) error {
	err := redis.RedisDb.HSet(session.sid, key, value)
	err = redis.RedisDb.Expire(session.sid, LifeTime)
	return err
}

func (session *session) SingleDel(key string) error {
	return redis.RedisDb.HDelete(session.sid, key)
}

func (session *session) Del() error {
	return redis.RedisDb.Delete(session.sid)
}

func genSid() string {
	strBuilder := strings.Builder{}
	strBuilder.Grow(SidLength)
	var str = "0123456789abcdefghijklmnopqrstuvwxyz"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < SidLength; i++ {
		strBuilder.WriteByte(str[r.Intn(len(str))])
	}
	return strBuilder.String()
}
