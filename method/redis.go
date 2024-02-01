package method

import (
	"github.com/gomodule/redigo/redis"
	"stock/biz/dal/sql"
	"stock/constant"
	"stock/util"
)

func DoGetRedisValue(key string) (errCode util.HttpCode, value interface{}) {
	if key == "" {
		errCode = util.HttpCode{
			Code:    constant.ERRDATALOSE,
			Message: "DoGetRedisValue 关键信息丢失",
			Data:    struct{}{},
		}
		return
	}

	value, err := redis.String(config.RedisConn.Get().Do("GET", key))
	if err != nil {
		errCode = util.HttpCode{
			Code:    constant.ERRDOREDIS,
			Message: "DoGetRedisValue 获取redis失败",
			Data:    struct{}{},
		}
		return
	}

	errCode = util.HttpCode{
		Code:    constant.ERRSUCCER,
		Message: "DoGetRedisValue 获取redis成功",
		Data:    struct{}{},
	}
	return
}

func DoDelRedisValue(key string) (errCode util.HttpCode) {
	if key == "" {
		errCode = util.HttpCode{
			Code:    constant.ERRDATALOSE,
			Message: "DoGetRedisValue 关键信息丢失",
			Data:    struct{}{},
		}
		return
	}

	_, err := redis.String(config.RedisConn.Get().Do("DEL", key))
	if err != nil {
		errCode = util.HttpCode{
			Code:    constant.ERRDOREDIS,
			Message: "DoDelRedisValue 获取redis失败",
			Data:    struct{}{},
		}
		return
	}

	errCode = util.HttpCode{
		Code:    constant.ERRSUCCER,
		Message: "DoGetRedisValue 获取redis成功",
		Data:    struct{}{},
	}
	return
}

func DoSetRedisValue(key string, value string, time int) (errCode util.HttpCode) {
	if key == "" || value == "" {
		errCode = util.HttpCode{
			Code:    constant.ERRDATALOSE,
			Message: "DoGetRedisValue 关键信息丢失",
			Data:    struct{}{},
		}
		return
	}

	if time == 0 {
		_, err := config.RedisConn.Get().Do("SET", key, value)
		if err != nil {
			errCode = util.HttpCode{
				Code:    constant.ERRDOREDIS,
				Message: "DoSetRedisValue 存储redis失败",
				Data:    struct{}{},
			}
		}
	}
	if time != 0 {
		_, err := config.RedisConn.Get().Do("SETEX", key, time, value)
		if err != nil {
			errCode = util.HttpCode{
				Code:    constant.ERRDOREDIS,
				Message: "DoSetRedisValue 存储redis失败",
				Data:    struct{}{},
			}
		}
	}
	errCode = util.HttpCode{
		Code:    constant.ERRSUCCER,
		Message: "DoGetRedisValue 存储redis成功",
		Data:    struct{}{},
	}
	return
}
