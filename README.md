




запуск docker контейнера Nginx:

    $ docker run -it -d -p 80:80 -v /home/alex/server://usr/share/nginx/html:ro --name server nginx:alpine
