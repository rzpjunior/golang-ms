FROM git.edenfarm.id:5050/image/nginx-print/master:latest

WORKDIR /var/www/html
COPY . .

RUN cp nginx/apps.conf /etc/nginx/sites-enabled/default
RUN composer install --no-dev --optimize-autoloader
RUN php artisan optimize:clear
RUN php artisan view:clear
RUN php artisan key:generate
RUN chmod -R 777 storage bootstrap/cache

CMD php-fpm7.2 -D && nginx -g "daemon off;"
