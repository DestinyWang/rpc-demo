package services

type IUserService interface {
	GetName(userId int64) string
	DelUser(userId int64) error
	PutUser(userId int64, username string) error
}

var data = make(map[int64]string)

type UserService struct {
}

func (this *UserService) GetName(userId int64) string {
	if username, ok := data[userId]; ok {
		return username
	}
	return "guest"
}

func (this *UserService) DelUser(userId int64) error {
	delete(data, userId)
	return nil
}

func (this *UserService) PutUser(userId int64, username string) error {
	data[userId] = username
	return nil
}
