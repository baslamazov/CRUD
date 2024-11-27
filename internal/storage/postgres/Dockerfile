# Используем последнюю стабильную версию PostgreSQL  
FROM postgres:15  

# Устанавливаем рабочую директорию  
WORKDIR /docker-entrypoint-initdb.d  

# Копируем SQL-скрипты в контейнер  
COPY sql-scripts/ /docker-entrypoint-initdb.d/

# Устанавливаем права на выполнение скриптов  
RUN chmod -R 755 /docker-entrypoint-initdb.d  

# Указываем переменные окружения для настройки базы данных  
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres

# Указываем порты, которые будет слушать контейнер
#EXPOSE 5432

# Контейнер будет автоматически запускать PostgreSQL  
CMD ["postgres"]  