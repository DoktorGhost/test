#!/bin/bash

# Ожидание запуска PostgreSQL (может потребоваться подождать, пока база данных будет доступна)
until psql "user=admin dbname=people_database password=admin sslmode=disable" -c '\q'; do
  >&2 echo "PostgreSQL is unavailable - sleeping"
  sleep 1
done

# Выполнение миграции с использованием Goose
goose -dir db/migrations postgres "user=admin dbname=people_database password=admin sslmode=disable" up

# Запуск вашего Go-приложения
./main
