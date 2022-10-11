
## Описание
  websocket chat:
  1. аутентификация с использованием username


  3. создание чата с пользователем или открыть существующий
  
  4. создание группы 


  Не с первого раза определяет email из URL!!!

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
        id serial not null unique primary key, 
        username VARCHAR(255) not null unique,
        email VARCHAR(255) not null unique,
        password VARCHAR(255) not null
      );
   <!-- chats INTEGER[] -->

 <!-- для того, чтобы определять, с каким пользователем уже есть чат, создадим в таблице users поле с id пользователями, с котороми создан чат
  1 user

    {
    "username": "Alex",
    "email": "komalex",
    "password": "qwerty"
    } -->

  <!-- изменение массива:

    $ UPDATE users SET chats[cardinality(chats) + 1] = 1 WHERE id = 1;

  поиск в массиве:

    $ SELECT id FROM users WHERE {id_user2} = ANY (chats) AND id = {id_user1}; -->

_____________________________________________________
  Cоздание таблиц с созданными чатами и с историей чата:
   при создании чата, создается его id и записываются id пользователей

    $ create table if not exists chats
      ( 
        id serial not null unique primary key, 
        id_user1 integer references users (id) not null,
        id_user2 integer references users (id) not null
      );  

  Запись при создании чата:

    $ INSERT INTO chats (id_user1, id_user2) VALUES (1, 2);

  Проверка на существование чата между пользователями:

    $ SELECT * FROM chats WHERE id_user1=$1 id_user2=$2;

если вернется 1 строка, то таблица уже создана, если 0 - то еще нет.



  Также создается таблица для хранения истории чата, в названии которой применяется id чата и хранятся сообщения

    $ create table if not exists history_chat{id_chat}
      (  
        date timestamp,
        username VARCHAR(255) references users (username) on delete cascade not null,
        message VARCHAR(255) not null
      );      


  Запись в таблицу истории чата:

    $ INSERT INTO history_chat{id_chat} (date, username, message) VALUES (TIMESTAMP '2004-10-19 10:23:54', 'Alex', 'Hello, Bob!');

______________________________________________________
Создание группового чата:
  таблица с группами

    $ create table if not exists groups
      ( 
        id serial not null unique primary key, 
        title VARCHAR(255) not null,
        admin integer references users (id) on delete cascade not null
      );  

таблица участников групп и чатов:

    $ create table if not exists user_group
      ( 
        id_user integer references users (id) on delete cascade,
        id_group integer references groups (id) on delete cascade
      );  

  Также создается таблица для хранения истории группы, в названии которой применяется id группы и хранятся сообщения

    $ create table if not exists history_group{id_group}
      (  
        date timestamp,
        username VARCHAR(255) references users (username) on delete cascade not null,
        message VARCHAR(255) not null
      );  