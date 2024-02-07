package method

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"google.golang.org/appengine/log"
	"stock/biz/dal/sql"
	"stock/constant"
	"stock/util"
)

func DoGetRedisValue(c *gin.Context, key string) (errCode util.HttpCode, value interface{}) {
	if key == "" {
		log.Errorf(c, "DoGetRedisValue 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return
	}

	value, err := redis.String(config.RedisConn.Get().Do("GET", key))
	if err != nil {
		log.Errorf(c, "DoGetRedisValue 获取redis失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOREDIS,
			Data: struct{}{},
		}
		return
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return
}

func DoDelRedisValue(c *gin.Context, key string) (errCode util.HttpCode) {
	if key == "" {
		log.Errorf(c, "DoGetRedisValue 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return
	}

	_, err := redis.String(config.RedisConn.Get().Do("DEL", key))
	if err != nil {
		log.Errorf(c, "DoDelRedisValue 获取redis失败 err%d", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOREDIS,
			Data: struct{}{},
		}
		return
	}

	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return
}

func DoSetRedisValue(c *gin.Context, key string, value string, time int) (errCode util.HttpCode) {
	if key == "" || value == "" {
		log.Errorf(c, "DoGetRedisValue 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
		return
	}

	if time == 0 {
		_, err := config.RedisConn.Get().Do("SET", key, value)
		if err != nil {
			log.Errorf(c, "DoSetRedisValue 存储redis失败 err%d", err)
			errCode = util.HttpCode{
				Code: constant.ERRDOREDIS,
				Data: struct{}{},
			}
		}
	}
	if time != 0 {
		_, err := config.RedisConn.Get().Do("SETEX", key, time, value)
		if err != nil {
			log.Errorf(c, "DoSetRedisValue 存储redis失败 err%d", err)
			errCode = util.HttpCode{
				Code: constant.ERRDOREDIS,
				Data: struct{}{},
			}
		}
	}
	errCode = util.HttpCode{
		Code: constant.ERRSUCCER,
		Data: struct{}{},
	}
	return
}
