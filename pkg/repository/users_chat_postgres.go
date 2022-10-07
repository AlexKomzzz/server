package repository

import (
	"errors"
	"fmt"

	chat "github.com/AlexKomzzz/server"
)

// создание таблицы для хранения истории чата с другим пользователем
// в названии таблицы первый id пользователя ВСЕГДА меньше второго
func (r *Repository) CreateChat(idUser1, idUser2 int) error {
	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	query := fmt.Sprintf(`create table if not exists chat%d%d
				( 
					id serial not null unique, 
					date timestamp,
					username VARCHAR(255) references users  on delete cascade not null,
					message VARCHAR(255) not null
					);`, idUser1, idUser2)

	// создание таблицы. Если уже создана то пропуск
	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("error: 'exec' from CreateChat (repos): %v", err)
	}

	return nil
}

// определение, создан ли чат между этими пользователями
func (r *Repository) GetChatsByUsers(idUser1, idUser2 int) (bool, error) {
	var result bool

	query := "SELECT id FROM users WHERE $1 = ANY (chats) AND id = $2"
	res, err := r.db.Exec(query, idUser2, idUser1)
	if err != nil {
		return false, err
	}

	// выведем кол-во возвращаемых строк
	numRows, err := res.RowsAffected()
	if err != nil {
		return false, err
	} else if numRows > 1 {
		return false, errors.New("error: больше одного чата с этим пользователем")
	} else if numRows == 1 {
		result = true
	}

	return result, nil
}

// добавление id пользователей в колонку chats каждого из пользователей
func (r *Repository) SetChatsByUser(idUser1, idUser2 int) error {

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
}

// выгрузка истории чата
func (r *Repository) GetHistoryChat(historyChat *[]chat.Message, idUser1, idUser2 int) ([]chat.Message, error) {
	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	//historyChat := make([]*chat.Message, 0)

	query := fmt.Sprintf("SELECT * FROM chat%d%d", idUser1, idUser2)
	err := r.db.Select(historyChat, query)
	if err != nil {
		return nil, fmt.Errorf("error: 'select' from GetHistoryChat (repos): %v", err)
	}
	return *historyChat, nil
}

// добавление записи в чат
func (r *Repository) WriteInChat(msg *chat.Message, idUser1, idUser2 int) error {

	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	query := fmt.Sprintf("INSERT INTO chat%d%d (date, username, message) VALUES (TIMESTAMP '$1', $2, $3)", idUser1, idUser2)
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
