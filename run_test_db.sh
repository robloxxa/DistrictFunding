docker run -d --name postgrestest -p 5432:5432 -e POSTGRES_PASSWORD=test -e POSTGRES_USER=test -e POSTGRES_DB=test -v ./db:/docker-entrypoint-initdb.d postgres:alpine


docker run -d --name pgadmin-dev -e PGADMIN_DEFAULT_EMAIL=robloxxa@yandex.ru -e PGADMIN_DEFAULT_PASSWORD=test -p 82:80 dpage/pgadmin4:latest