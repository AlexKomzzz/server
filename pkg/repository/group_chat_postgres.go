package repository

import (
	"fmt"
	"log"

	chat "github.com/AlexKomzzz/server"
)

// создание группового чата
// возвращает id созданной группы
func (r *Repository) CreateGroup(title string, idAdmin int) (int, error) {

	var idGroup int

	// создание транзакции
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'Begin' from CreateGroup (repos): %v", err)
	}
	defer tx.Rollback()

	// добавление параметров новой группы в таблицу groups
	query := "INSERT INTO groups (title, admin) VALUES ($1, $2) RETURNING id"
	row := tx.QueryRow(query, title, idAdmin)
	if err := row.Scan(&idGroup); err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'Scan' from CreateGroup (repos): %v", err)
	}
	log.Printf("добавление группы №%d в таблицу groups\n", idGroup)

	//создание таблицы для хранения истории чата группы
	queryCreate := fmt.Sprintf(`create table if not exists history_group%d
							(  
	  							date timestamp,
	  							username VARCHAR(255) references users (username) on delete cascade not null,
	 							message VARCHAR(255) not null
							)`, idGroup)

	// создание таблицы. Если уже создана то пропуск
	_, err = tx.Exec(queryCreate)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error при создании таблицы history_group: func 'Exec CREATE' from CreateGroup (repos): %v", err)
	}
	log.Printf("создание таблицы history_group№%d\n", idGroup)

	// добавление записи в таблицу user_group
	queryInsert := "INSERT INTO user_group (id_user, id_group) VALUES ($1, $2)"
	_, err = tx.Exec(queryInsert, idAdmin, idGroup)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error при собавление записи в таблицу user_group: func 'Exec INSERT' from CreateGroup (repos): %v", err)
	}
	log.Println("добавление записи в таблицу user_group")

	// коммит транзакции
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'Commit' from UpdateChats (repos): %v", err)
	}

	return idGroup, nil
}

// выгрузка истории группового чата
func (r *Repository) GetHistoryGroup(idGroup int) ([]*chat.Message, error) {

	historyChat := make([]*chat.Message, 0)

	query := fmt.Sprintf("SELECT * FROM history_group%d", idGroup)
	err := r.db.Select(&historyChat, query)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return nil, fmt.Errorf("error: 'select' from GetHistoryChat (repos): %v", err)
		}
	}

	return historyChat, nil
}

// добавление записи в групповой чат
func (r *Repository) WriteInGroupChat(msg *chat.Message, idGroup int) error {

	query := fmt.Sprintf("INSERT INTO history_group%d (date, username, message) VALUES ($1, $2, $3)", idGroup)
	// ничего не возвращаем, используем exec
	_, err := r.db.Exec(query, msg.Date, msg.Username, msg.Body)
	if err != nil {
		return fmt.Errorf("error: 'exec' from WriteInGroupChat (repos): %v", err)
	}
	// numRows, err := res.RowsAffected()
	// if err != nil {
	// 	return fmt.Errorf("error: 'RowsAffected' from WriteInChat (repos): %v", err)
	// } else if numRows == 0 {
	// 	return fmt.Errorf("error: сообщение не записалось в БД: ")
	// }

	return nil
}
