// api
package api

import (
	"fmt"
	"net/http"

	"time"

	"lottery-server/constants"
	"lottery-server/model"
	"lottery-server/response"
	"lottery-server/util"

	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func NewTopic(c *gin.Context) {

	var response response.BaseReponse

	userId := c.PostForm("userId")
	userToken := c.PostForm("userToken")
	topicContent := c.PostForm("TopicContent")
	topicTitle := c.PostForm("TopicTitle")
	topicTag := c.PostForm("TopicTag")

	if userId == "" || userToken == "" || topicTitle == "" {
		response.Status = constants.ApiCodeFail
		response.Message = "参数不正确"

		c.JSON(http.StatusOK, response)

		return
	}

	db, err := util.OpenLotteryDb()
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	rows, err := db.Query("SELECT * FROM user WHERE UserId=? AND UserToken=?", userId, userToken)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	if rows.Next() == false {
		response.Status = constants.ApiCodeFail
		response.Message = "用户口令错误"

		c.JSON(http.StatusOK, response)
		db.Close()
		return
	}

	stmt, err := db.Prepare("INSERT topic SET TopicId=?,UserId=?,TopicTitle=?,TopicContent=?,TopicTag=?,LastModifyTime=?")
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	topicId := util.UUID()
	lastModifyTime := time.Now().Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(topicId, userId, topicTitle, topicContent, topicTag, lastModifyTime)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	db.Close()

	response.Status = constants.ApiCodeSuccess
	response.Message = "success"

	c.JSON(http.StatusOK, response)
}

func UpdateTopic(c *gin.Context) {

	var response response.BaseReponse

	userId := c.PostForm("userId")
	userToken := c.PostForm("userToken")
	topicContent := c.PostForm("TopicContent")
	topicTitle := c.PostForm("TopicTitle")
	topicTag := c.PostForm("TopicTag")
	topicId := c.PostForm("TopicId")

	if userId == "" || userToken == "" || topicTitle == "" || topicId == "" {
		response.Status = constants.ApiCodeFail
		response.Message = "参数不正确"

		c.JSON(http.StatusOK, response)

		return
	}

	db, err := util.OpenLotteryDb()
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	rows, err := db.Query("SELECT * FROM user WHERE UserId=? AND UserToken=?", userId, userToken)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	if rows.Next() == false {
		response.Status = constants.ApiCodeFail
		response.Message = "用户口令错误"

		c.JSON(http.StatusOK, response)
		db.Close()
		return
	}

	rows, err = db.Query("SELECT * FROM topic WHERE UserId=? AND TopicId=?", userId, topicId)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	if rows.Next() == false {
		response.Status = constants.ApiCodeFail
		response.Message = "无操作权限或帖子不存在"

		c.JSON(http.StatusOK, response)
		db.Close()
		return
	}

	stmt, err := db.Prepare("UPDATE topic SET TopicTitle=?,TopicContent=?,TopicTag=?,LastModifyTime=? WHERE TopicId=?")
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	lastModifyTime := time.Now().Format("2006-01-02 15:04:05")

	_, err = stmt.Exec(topicTitle, topicContent, topicTag, lastModifyTime, topicId)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	db.Close()

	response.Status = constants.ApiCodeSuccess
	response.Message = "success"

	c.JSON(http.StatusOK, response)
}

func QueryTopic(c *gin.Context) {

	var response response.QueryTopicListReponse

	userId := c.Query("userId")
	topicTitle := c.Query("topicTitle")
	topicTag := c.Query("topicTag")

	db, err := util.OpenLotteryDb()
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	var conditions []string
	var values []interface{}
	var rows *sql.Rows
	var basesql = "SELECT * FROM topic"

	if userId != "" {
		conditions = append(conditions, " UserId=?")
		values = append(values, userId)

	}
	if topicTitle != "" {
		conditions = append(conditions, " TopicTitle LIKE ?")
		values = append(values, "%"+topicTitle+"%")

	}

	if topicTag != "" {
		conditions = append(conditions, " TopicTag=?")
		values = append(values, topicTag)
	}

	for index, val := range conditions {

		if index != 0 {

			basesql = basesql + " AND"
		} else {
			basesql = basesql + " WHERE"
		}
		basesql = basesql + val
	}

	fmt.Println("basesql:" + basesql)

	rows, err = db.Query(basesql, values...)
	if util.CheckHttpErr(err, c) == false {
		db.Close()
		return
	}

	for rows.Next() {
		var topic model.Topic

		err = rows.Scan(&topic.TopicId, &topic.UserId, &topic.TopicTitle, &topic.TopicContent, &topic.TopicTag, &topic.LastModifyTime)

		if util.CheckHttpErr(err, c) == false {
			db.Close()
			return
		}
		topic.UserName = "lic"

		response.Data = append(response.Data, topic)
	}

	db.Close()

	response.Status = constants.ApiCodeSuccess
	response.Message = "success"

	c.JSON(http.StatusOK, response)
}
