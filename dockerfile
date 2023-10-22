# Используем официальный образ PostgreSQL
FROM postgres:latest

# Устанавливаем необходимые переменные окружения для PostgreSQL
ENV POSTGRES_USER admin
ENV POSTGRES_PASSWORD admin
ENV POSTGRES_DB people_database

# Экспонируем порт PostgreSQL (по умолчанию 5432)
EXPOSE 5432

