package repository

import (
	"errors"
	"fmt"

	chat "github.com/AlexKomzzz/server"
)

// определение, был ли чат между этими пользователями создан ранее
func (r *Repository) getChatsByUsers(idUser, idNewUser int) (bool, error) {
	var result bool

	query := "SELECT id FROM users WHERE $1 = ANY (chats) AND id = $2"
	res, err := r.db.Exec(query, idNewUser, idUser)
	if err != nil {
		return false, err
	}
	numRows, err := res.RowsAffected()
	// numRows - кол-во возвращаемых строк
	if err != nil {
		return false, err
	} else if numRows > 1 {
		return false, errors.New("error: больше одного чата с этим пользователем")
	}
	// else if numRows == 0 {
	// 	err := r.setChatsByUser(idUser, idNewUser)
	// 	if err != nil {
	// 		return false, err
	// 	}

	// 	result = true
	// }

	return result, nil
}

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
					user_1 integer references users (id) not null,
					user_2 integer references users (id) not null,
					date timestamp,
					username VARCHAR(255) references users  on delete cascade not null,
					message VARCHAR(255) not null
				);`, idUser1, idUser2)

	_, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

// выгрузка истории чата
func (r *Repository) GetChat(idUser1, idUser2 int) ([]*chat.Message, error) {
	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	historyChat := make([]*chat.Message, 0)

	query := fmt.Sprintf("SELECT (id, date, username, message) FROM chat%d%d", idUser1, idUser2)
	err := r.db.Select(&historyChat, query)
	if err != nil {
		return nil, fmt.Errorf("error: 'select' from GetChat (repos): %v", err)
	}
	return historyChat, nil
}

// добавление записи в чат
func (r *Repository) WriteInChat(idUser1, idUser2 int, msg *chat.Message) error {

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

// добавление id нового пользователя в колонку chats текущего пользователя
func (r *Repository) setChatsByUser(idUser, idNewUser int) error {
	query := "UPDATE users SET chats[cardinality(chats) + 1] = $1 WHERE id = $2"
	res, err := r.db.Exec(query, idNewUser, idUser)
	if err != nil {
		return err
	}
	numRows, err := res.RowsAffected()
	if err != nil {
		return err
	} else if numRows == 0 {
		return errors.New("warning: пользователю не добавился id чата")
	}

	// второму пользователю также добавим id
	res, err = r.db.Exec(query, idUser, idNewUser)
	if err != nil {
		return err
	}
	numRows, err = res.RowsAffected()
	if err != nil {
		return err
	} else if numRows == 0 {
		return errors.New("warning: пользователю не добавился id чата")
	}

	return nil
}
