package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"stock/constant"
)

var store = sessions.NewCookieStore([]byte(constant.SessionEncipher))
var sessionName = constant.SessionName

func GetSession(c *gin.Context) map[interface{}]interface{} {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session:%+v\n", session.Values)
	return session.Values
}

// SetSession 创建session
func SetSession(c *gin.Context, userName string, userId string, ip string) error {
	session, _ := store.Get(c.Request, sessionName)
	session.Values["userName"] = userName
	session.Values["userId"] = userId
	session.Values["ip.json"] = ip
	return session.Save(c.Request, c.Writer)
}

func FlushSession(c *gin.Context) error {
	session, _ := store.Get(c.Request, sessionName)
	fmt.Printf("session : %+v\n", session.Values)
	session.Values["userName"] = ""
	session.Values["ip.json"] = ""
	session.Values["userId"] = 0
	return session.Save(c.Request, c.Writer)
}
