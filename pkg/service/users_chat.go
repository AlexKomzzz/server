package service

import chat "github.com/AlexKomzzz/server"

// Создание чата с клиентом
func (service *Service) GetChat(historyChat []*chat.Message, idUser1 int, emailUser2 string) ([]*chat.Message, error) {

	// определение id второго пользователя по email
	idUser2, err := service.repos.GetUserByEmail(emailUser2)
	if err != nil {
		return nil, err
	}

	// определение, создан ли чат между этими пользователями
	if ok, err := service.repos.GetChatsByUsers(idUser1, idUser2); !ok {
		if err != nil {
			return nil, err
		}

		// если чата нет, то создадим его
		err := service.repos.CreateChat(idUser1, idUser2)
		if err != nil {
			return nil, err
		}

		// и добавим пользователей друг другу в поле chats
		err = service.repos.SetChatsByUser(idUser1, idUser2)
		if err != nil {
			return nil, err
		}
	}

	// получение истории чата пользователей
	service.repos.GetHistoryChat(&historyChat, idUser1, idUser2)

	return historyChat, nil
}
