package service

import chat "github.com/AlexKomzzz/server"

// создание группового чата
func (service *Service) CreateGroup(title string, idUser int) (int, error) {

	idGroup, err := service.repos.CreateGroup(title, idUser)
	if err != nil {
		return -1, err
	}

	return idGroup, nil
}

// получение истории группового чата
func (service *Service) GetGroup(idGroup int) ([]*chat.Message, error) {

	// получение истории чата пользователей
	return service.repos.GetHistoryGroup(idGroup)
}

// сохранение нового сообщения в чат
func (service *Service) WriteInGroup(msg *chat.Message, idGroup int) error {
	err := service.repos.WriteInGroupChat(msg, idGroup)
	if err != nil {
		return err
	}

	return nil
}

// получение id всех созданных групповых чатов
func (service *Service) GetIdGroups() ([]int, error) {

	// получение истории чата пользователей
	return service.repos.GetIdGroups()
}
