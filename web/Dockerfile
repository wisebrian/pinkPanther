FROM ubuntu
RUN apt-get update
RUN apt-get install nginx -y
COPY index.html /var/www/html
EXPOSE 8082
CMD ["nginx","-g","daemon off;"]
RUN sed -i 's/80/8082/1' /etc/nginx/sites-enabled/default