package services

import (
	"fmt"
	"github.com/DestinyWang/gokit-test/util"
)

type IUserService interface {
	GetName(userId int64) string
	DelUser(userId int64) error
	PutUser(userId int64, username string) error
}

//var data = map[int64]string{
//	101: "destiny",
//}

type UserService struct {
}

func (this *UserService) GetName(userId int64) string {
	//if username, ok := data[userId]; ok {
	//	return username
	//}
	//return "guest"
	switch userId {
	case 101:
		return fmt.Sprintf("Destiny:%d", util.ServicePort)
	case 102:
		return "Freedom"
	default:
		return "Guest"
	}
}

func (this *UserService) DelUser(userId int64) error {
	//delete(data, userId)
	return nil
}

func (this *UserService) PutUser(userId int64, username string) error {
	//data[userId] = username
	return nil
}
