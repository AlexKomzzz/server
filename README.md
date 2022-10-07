
## Описание
  websocket chat:
  1. аутентификация с использованием username
  hendler на изменение username

  3. создание чата с пользователем или открыть существующий
  для того, чтобы определять, с каким пользователем уже есть чат, создадим в таблице users поле с id пользователями, с котороми создан чат

## Docker

создание тома:

docker volume create --driver local \
  --opt type=none \
  --opt device=/var/www/alexkomzzz.ml \
  --opt o=bind \
  domen

  --mount type=volume,source=home/alex/server/nginx/logs,destination=/usr/share/nginx/html,ro

запуск docker контейнера Nginx:

    $ docker run -d --network server-net --ip 172.16.0.4 --name nginx -p 80:80 -v domen:/var/www/alexkomzzz.ml -v srv-log:/var/log/nginx -v srv-conf:/etc/nginx/conf.d nginx

docker exec nginx ls /var/www/alexkomzzz.ml


docker run --name server --mount type=bind,source=/home/alex/server/nginx/test,destination=/usr/share/nginx/html,ro -d nginx
-v /home/alex/server/nginx/log:/var/log/nginx/


ln -s /etc/nginx/sites-available/alexkomzzz.ml.conf /etc/nginx/sites-enabled/


/usr/share/nginx/html


 БД
docker exec -it db /bin/bash
psql -U postgres

## Структура БД

1 таблица: Users
  id, username, email, password_hash

2 таблица: history[id_1_user][id_2_user] - для каждого нового диалога создается новая таблица
  data, times, username, message

    $ create table if not exists users
      ( 
        id serial not null unique, 
        username VARCHAR(255) not null primary key,
        email VARCHAR(255) not null unique,
        password VARCHAR(255) not null,
        chats INTEGER[]
      );


1 user

    {
    "username": "Alex",
    "email": "komalex",
    "password": "qwerty"
    }

  изменение массива:

    $ UPDATE users SET chats[cardinality(chats) + 1] = 1 WHERE id = 1;

  поиск в массиве:

    $ SELECT id FROM users WHERE {id_user2} = ANY (chats) AND id = {id_user1};

  создание таблицы с историей чата:

    $ create table if not exists chat12
      ( 
        id serial not null unique, 
        date timestamp,
        username VARCHAR(255) references users  on delete cascade not null,
        message VARCHAR(255) not null
      );      

//         id_user integer references users (id) not null,

  Запись в таблицу чата:

    $ INSERT INTO chat12 (date, username, message) VALUES (TIMESTAMP '2004-10-19 10:23:54', 'Alex', 'Hello, Bob!');