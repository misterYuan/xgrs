/*
基于gin框架的session
*/
package xgrs

import (
	"errors"
	"strconv"

	"github.com/misterYuan/xrdb"
	"github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

var session_expire_err = errors.New("会话过期")

type ginS struct {
	ctx       *gin.Context
	sessionId string
}

func GinSStart(ctx *gin.Context) *ginS {
	return &ginS{ctx, ""}
}

type sconfig struct {
	cookieName string
	maxAge     int
	secure     bool
	httpOnly   bool
}

var sf *sconfig

func SetSConfig(cookieName string, maxAge int, secure, httpOnly bool) {
	sf = &sconfig{cookieName, maxAge, secure, httpOnly}
}

func getSessionId() string {
	return uuid.Must(uuid.NewV4()).String()
}

func (g *ginS) Set(key, value string) *ginS {
	if g.sessionId == "" {
		var err error
		if g.sessionId, err = g.ctx.Cookie(sf.cookieName); err != nil {
			g.sessionId = getSessionId()
		}
	}
	g.ctx.SetCookie(sf.cookieName, g.sessionId, sf.maxAge, "/", "", sf.secure, sf.httpOnly)
	if err := xrdb.HMSet(g.sessionId, key, value, strconv.Itoa(sf.maxAge)); err != nil {
		panic(err.Error())
	}
	return g
}

func (g *ginS) Get(key string) (string, error) {
	si, err := g.ctx.Cookie(sf.cookieName)
	if err != nil {
		return "", err
	}
	value, ok := xrdb.HMGet(si, key)
	if !ok {
		return "", session_expire_err
	}
	return value, nil
}
