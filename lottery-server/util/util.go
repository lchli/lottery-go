// util
package util

import (
	"database/sql"

	"lottery-server/constants"

	"lottery-server/response"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

var create_user_table_sql = "CREATE TABLE IF NOT EXISTS `user` (" +
	"`UserId` VARCHAR(200) NOT NULL PRIMARY KEY," +
	"`UserName` VARCHAR(100) NOT NULL," +
	"`UserPwd` VARCHAR(200) NOT NULL," +
	"`UserToken` VARCHAR(200)" +
	");"

var create_topic_table_sql = "CREATE TABLE IF NOT EXISTS `topic` (" +
	"`TopicId` VARCHAR(200) NOT NULL PRIMARY KEY," +
	"`UserId` VARCHAR(100) NOT NULL," +
	"`TopicTitle` VARCHAR(200) NOT NULL," +
	"`TopicContent` TEXT," +
	"`TopicTag` VARCHAR(100)," +
	"`LastModifyTime` VARCHAR(100) NOT NULL" +
	");"

func OpenLotteryDb() (*sql.DB, error) {

	db, err := sql.Open("mysql", "root@/biz")

	if err != nil {
		return db, err
	}

	db, err = createTable(create_user_table_sql, db)
	if err != nil {
		return db, err
	}

	db, err = createTable(create_topic_table_sql, db)
	if err != nil {
		return db, err
	}

	return db, err
}

func createTable(sql string, db *sql.DB) (*sql.DB, error) {
	stmt, err := db.Prepare(sql)

	if err != nil {
		return db, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return db, err
	}
	return db, err
}

func CheckHttpErr(err error, c *gin.Context) bool {
	if err != nil {

		c.JSON(http.StatusOK, response.BaseReponse{constants.ApiCodeFail, err.Error()})
		return false
	}

	return true
}

func UUID() string {
	return uuid.NewV4().String()
}
