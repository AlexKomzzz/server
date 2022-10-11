package repository

import "fmt"

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

	// добавление записи в таблицу user_group
	queryInsert := "INSERT INTO user_group (id_user, id_group) VALUES ($1, $2)"
	_, err = tx.Exec(queryInsert, idAdmin, idGroup)
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error при собавление записи в таблицу user_group: func 'Exec INSERT' from CreateGroup (repos): %v", err)
	}

	// коммит транзакции
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, fmt.Errorf("error: 'Commit' from UpdateChats (repos): %v", err)
	}

	return idGroup, nil
}
