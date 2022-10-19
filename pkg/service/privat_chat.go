package service

import (
	chat "github.com/AlexKomzzz/server"
)

// Создание чата с клиентом(если его нет)
/*func (service *Service) CreatePrivChat(idUser1, idUser2 int) (int, error) {

	// создание приватного чата
	idChat, err := service.repos.CreatePrivChat(idUser1, idUser2)
	if err != nil {
		return -1, err
	}

	return idChat, nil
}*/

// Создание чата с клиентом(если его нет) и проверка, создан ли чат между пользователями
func (service *Service) CreateAndGetIdPrivChat(idUser1, idUser2 int) (int, error) {

	// создание приватного чата
	idChat, err := service.repos.GetIdPrivChat(idUser1, idUser2)
	if err != nil {
		return -1, err
	}

	return idChat, nil
}

// получение истории чата
func (service *Service) GetPrivChat(idUser1, idUser2 int) ([]*chat.Message, error) {

	idChat, err := service.repos.GetIdPrivChat(idUser1, idUser2)
	if err != nil {
		return nil, err
	}

	// получение истории чата пользователей
	return service.repos.GetHistoryChat(idChat)
}

// сохранение нового сообщения в чат
func (service *Service) WriteInPrivChat(msg *chat.Message, idChat int) error {

	// сохранение нового сообщения в БД
	err := service.repos.WriteInPrivChat(msg, idChat)
	if err != nil {
		return err
	}

	return nil
}

// получение id всех созданных приватных чатов
func (service *Service) GetIdsPrivChats() ([]int, error) {

	// получение истории чата пользователей
	return service.repos.GetIdsPrivChats()
}
