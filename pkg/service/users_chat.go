package service

// для того, чтобы определять, с каким пользователем уже есть чат, создадим в таблице users поле с id пользователями, с котороми создан чат
func (service *Service) GetChat(id_user1 int, email_user2 string) error {

	// определение id второго пользователя по email
	id_user2, err := service.repos.GetUserByEmail(email_user2)
	if err != nil {
		return err
	}

	// создание таблицы чата в бд, если ее нет
	err = service.repos.CreateChat(id_user1, id_user2)
	if err != nil {
		return err
	}

	return nil
}
