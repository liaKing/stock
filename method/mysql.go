package method

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"stock/biz/dal/sql"
	"stock/biz/model"
	"stock/constant"
	"stock/util"
)

// DoFindMySQLUser 查找数据库用户
func DoFindMySQLUser(c *gin.Context, UserName string) (errCode util.HttpCode, user *model.User) {
	user = &model.User{}
	query := `SELECT * FROM "user" WHERE user_name = $1`
	err := config.DB.Raw(query, UserName).First(user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("[ERROR] DoFindMySQLUser 操作mysql失败 err%v", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
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

// DoCreateMySQLUser 添加用户信息
func DoCreateMySQLUser(c *gin.Context, user *model.User) (errCode util.HttpCode) {
	if user == nil {
		log.Printf("[ERROR] DoCreateMySQLUser 关键信息丢失")
		errCode = util.HttpCode{
			Code: constant.ERRDATALOSE,
			Data: struct{}{},
		}
	}
	err := config.DB.Table("user").Create(user).Error
	if err != nil {
		log.Printf("[ERROR] DoCreateMySQLUser 操作mysql失败")
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
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

//通过id获取用户信息

func GetUserById(c *gin.Context, id string) (errCode util.HttpCode, user *model.User) {
	user = &model.User{}
	query := `SELECT * FROM "user" WHERE user_id = $1`
	err := config.DB.Raw(query, id).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.HttpCode{Code: constant.ERRISNOTEXIT, Data: struct{}{}}, nil
		}
		return util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}, nil
	}
	return util.HttpCode{Code: constant.ErrSuccer, Data: struct{}{}}, user
}

// DoUpdataMySQLUser 修改数据库用户
func DoUpdataMySQLUser(c *gin.Context, userId string, deletionReason string) (errCode util.HttpCode) {
	query := `UPDATE "user" SET del_flg = $1, deletion_reason = $2 WHERE user_id = $3`
	err := config.DB.Exec(query, 1, deletionReason, userId).Error
	if err != nil {
		log.Printf("[ERROR] DoUpdataMySQLUser 操作mysql失败 err%v", err)
		errCode = util.HttpCode{
			Code: constant.ERRDOMYSQL,
			Data: struct{}{},
		}
		return
	}
	errCode = util.HttpCode{
		Code: constant.ErrSuccer,
		Data: struct{}{},
	}

	return
}

// DoUpdataMySQLUserLuck 修改数据库用户幸运值（luck 传 int64 避免 PostgreSQL BIGINT 参数类型歧义）
func DoUpdataMySQLUserLuck(c *gin.Context, userId string, luck uint64) (errCode util.HttpCode) {
	query := `UPDATE "user" SET luck = $1 WHERE user_id = $2`
	result := config.DB.Exec(query, int64(luck), userId)
	if result.Error != nil {
		log.Printf("[ERROR] DoUpdataMySQLUserLuck 操作失败 err:%v", result.Error)
		return util.HttpCode{Code: constant.ERRDOMYSQL, Data: struct{}{}}
	}
	if result.RowsAffected == 0 {
		log.Printf("[ERROR] DoUpdataMySQLUserLuck 未命中任何行 userId:%s luck:%d", userId, luck)
		return util.HttpCode{Code: constant.ERRISNOTEXIT, Data: struct{}{}}
	}
	return util.HttpCode{Code: constant.ErrSuccer, Data: struct{}{}}
}
