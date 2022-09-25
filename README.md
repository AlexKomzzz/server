
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