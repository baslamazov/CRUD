# Используем последнюю стабильную версию PostgreSQL  
FROM postgres:15  

WORKDIR /docker-entrypoint-initdb.d  

COPY sql-scripts/ /docker-entrypoint-initdb.d/

RUN chmod -R 755 /docker-entrypoint-initdb.d  

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=postgres

CMD ["postgres"]  