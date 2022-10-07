package repository

import (
	"errors"
	"fmt"
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
func (r *Repository) createChat(idUser1, idUser2 int) error {
	// определяем меньший id пользователя
	if idUser1 > idUser2 {
		idUser1, idUser2 = idUser2, idUser1
	}

	query := fmt.Sprintf(`create table if not exists chat%d%d
				( 
					id serial not null unique, 
					user_1 integer references users (id) not null,
					user_2 integer references users (id) not null,
					data timestamp,
					username VARCHAR(255) references users  on delete cascade not null,
					message VARCHAR(255) not null
				);`, idUser1, idUser2)

	_, err := r.db.Exec(query)
	if err != nil {
		return err
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
