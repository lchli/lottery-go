// model
package model

type User struct {
	UserId    string
	UserName  string
	UserPwd   string
	UserToken string
}

type Topic struct {
	UserId         string
	TopicId        string
	UserName       string
	TopicTitle     string
	TopicContent   string
	TopicTag       string
	LastModifyTime string
}
