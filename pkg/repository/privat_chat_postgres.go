package repository

import (
	"fmt"

	chat "github.com/AlexKomzzz/server"
)

// проверка, создан ли чат между этими пользователями
func (r *Repository) GetIdPrivChat(idUser1, idUser2 int) (int, error) {

	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	var idChat int

	query := "SELECT id FROM chats WHERE id_user1=$1 AND id_user2=$2"

	row := r.db.QueryRow(query, idUser1, idUser2)
	err := row.Scan(&idChat)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {

			// значит чат не был создан
			idChat, err = r.CreatePrivChat(idUser1, idUser2)
			if err != nil {
				return -1, fmt.Errorf("%s", err)
			}

		} else {
			return -1, fmt.Errorf("error: 'exec' from CreateChatTwoUser (repos): %v", err)
		}
	}

	return idChat, nil
}

// первый шаг при создании чата между пользователями.
// Добавление записи в таблицу chats
// idUser1 < idUser2
func (r *Repository) CreatePrivChat(idUser1, idUser2 int) (int, error) {

	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	var idChat int
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'Begin' from UpdateChats (repos): %v", err)
	}
	defer tx.Rollback()

	query := "INSERT INTO chats (id_user1, id_user2) VALUES ($1, $2) RETURNING id"

	row := tx.QueryRow(query, idUser1, idUser2)
	err = row.Scan(&idChat)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: ошибка при возврате id чата в UpdateChats (repos): %v", err)
	}

	//создание таблицы для хранения истории чата с другим пользователем
	query = fmt.Sprintf("create table if not exists history_chat%d	(date timestamp, username VARCHAR(255) references users (username) on delete cascade not null,	message VARCHAR(255) not null)", idChat)

	// создание таблицы. Если уже создана то пропуск
	_, err = tx.Exec(query)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'exec' from UpdateChats (repos): %v", err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'Commit' from UpdateChats (repos): %v", err)
	}

	return idChat, nil

}

// Второй шаг при создании чата между пользователями.
// Создание таблицы для хранения истории чата с другим пользователем
func (r *Repository) CreateHistoryChat(idChat int) error {

	query := fmt.Sprintf("create table if not exists history_chat%d	(date timestamp, username VARCHAR(255) references users (username) on delete cascade not null,	message VARCHAR(255) not null)", idChat)

	// создание таблицы. Если уже создана то пропуск
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error: 'exec' from CreateChat (repos): %v", err)
	}

	return nil
}

// добавление id пользователей в колонку chats каждого из пользователей
/*func (r *Repository) SetChatsByUser(idUser1, idUser2 int) error {

	// создадим транзакцию
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error: 'db.Begin' from SetChatsByUser (repos): %v", err)
	}

	// первому пользователю в поле chats добавим id второго пользователя
	query := "UPDATE users SET chats[cardinality(chats) + 1] = $1 WHERE id = $2"
	res, err := tx.Exec(query, idUser2, idUser1)
	if err != nil {
		tx.Rollback()
		return err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
		// если возвращено 0 строк, значит добавление не произошло
	} else if numRows == 0 {
		tx.Rollback()
		return errors.New("warning: пользователю не добавился id чата")
	}

	// второму пользователю также добавим id
	res, err = tx.Exec(query, idUser1, idUser2)
	if err != nil {
		tx.Rollback()
		return err
	}
	numRows, err = res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	} else if numRows == 0 {
		tx.Rollback()
		return errors.New("warning: пользователю не добавился id чата")
	}

	return tx.Commit()
}*/

// выгрузка истории чата
func (r *Repository) GetHistoryChat(idChat int) ([]*chat.Message, error) {

	historyChat := make([]*chat.Message, 0)

	query := fmt.Sprintf("SELECT * FROM history_chat%d", idChat)
	err := r.db.Select(&historyChat, query)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, fmt.Errorf("error: 'select' from GetHistoryChat (repos): %v", err)
		}
	}

	return historyChat, nil
}

// добавление записи в чат
func (r *Repository) WriteInPrivChat(msg *chat.Message, idChat int) error {

	query := fmt.Sprintf("INSERT INTO history_chat%d (date, username, message) VALUES ($1, $2, $3)", idChat)
	// ничего не возвращаем, используем exec
	res, err := r.db.Exec(query, msg.Date, msg.Username, msg.Body)
	if err != nil {
		return fmt.Errorf("error: 'exec' from WriteInChat (repos): %v", err)
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error: 'RowsAffected' from WriteInChat (repos): %v", err)
	} else if numRows == 0 {
		return fmt.Errorf("error: сообщение не записалось в БД: ")
	}

	return nil
}
