// response
package response

import "lottery-server/model"

type BaseReponse struct {
	Status  int
	Message string
}

type RegisterReponse struct {
	BaseReponse
	Data model.User
}

type QueryTopicListReponse struct {
	BaseReponse
	Data []model.Topic
}
