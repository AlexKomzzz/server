package service

// создание группового чата
func (service *Service) CreateGroup(title string, idUser int) (int, error) {

	idGroup, err := service.repos.CreateGroup(title, idUser)
	if err != nil {
		return -1, err
	}

	return idGroup, nil
}

// подключение к групповому чату, получение истории сообщений
// func (service *Service) GetGroup(idgroup int, emailUser2 string) ([]*chat.Message, error) {

// }
