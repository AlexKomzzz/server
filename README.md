
docker volume create --driver local \
  --opt type=none \
  --opt device=/home/alex/server/nginx/test \
  --opt o=volume \
  test


создание тома:

docker volume create --driver local \
  --opt type=none \
  --opt device=/var/www/alexkomzzz.ml \
  --opt o=bind \
  domen

docker volume create --driver local \
  --opt type=none \
  --opt device=/var/www/apps \
  --opt o=bind \
  srv-conf

docker volume create --driver local \
  --opt type=none \
  --opt device=/var/log/nginx \
  --opt o=bind \
  srv-log

  --mount type=volume,source=home/alex/server/nginx/logs,destination=/usr/share/nginx/html,ro

запуск docker контейнера Nginx:

    $ docker run --name server -v test:/usr/share/nginx/html:ro -dp 80:80 nginx
    --network server-net --ip 172.16.0.4


docker run -d --network server-net --ip 172.16.0.4 --name nginx -p 80:80 -v domen:/var/www/alexkomzzz.ml -v srv-log:/var/log/nginx -v srv-conf:/etc/nginx/conf.d nginx

docker exec nginx ls /var/www/alexkomzzz.ml


docker run --name server --mount type=bind,source=/home/alex/server/nginx/test,destination=/usr/share/nginx/html,ro -d nginx
-v /home/alex/server/nginx/log:/var/log/nginx/

docker run --name server -v /home/alex/server/nginx/logs:/usr/share/nginx/html:ro -dp 80:80 nginx
docker exec server ls /usr/share/nginx/html



docker volume create --driver local \
  --opt type=none \
  --opt device=/home/alex/tom \
  --opt o=bind \
  srv






docker run -it -d --network server-net --ip 172.16.0.2  --name api srv


sudo cp -R /home/alex/server /var/apps/alexkomzzz.ml/volumes/etc/

ln -s /etc/nginx/sites-available/alexkomzzz.ml.conf /etc/nginx/sites-enabled/
