package service

import chat "github.com/AlexKomzzz/server"

// Создание чата с клиентом(если его нет), получение истории чата
func (service *Service) GetChat(idUser1 int, emailUser2 string) ([]chat.Message, error) {

	// определение id второго пользователя по email
	idUser2, err := service.repos.GetUserByEmail(emailUser2)
	if err != nil {
		return nil, err
	}

	// Cоздание чата между этими пользователями, если еще не создан
	idChat, err := service.repos.CreateChatTwoUser(idUser1, idUser2)
	if err != nil {
		return nil, err
	}

	// получение истории чата пользователей
	historyChat, err := service.repos.GetHistoryChat(idChat)
	if err != nil {
		return nil, err
	}

	return historyChat, nil
}

// сохранение нового сообщения в чат
func (service *Service) WriteInChat(msg *chat.Message, idUser1 int, emailUser2 string) error {

	// определение id второго пользователя по email
	idUser2, err := service.repos.GetUserByEmail(emailUser2)
	if err != nil {
		return err
	}

	// сохранение нового сообщения в БД
	err = service.repos.WriteInChat(msg, idUser1, idUser2)
	if err != nil {
		return err
	}

	return nil
}
