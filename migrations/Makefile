include .env

migrate-up:
	migrate -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_DATABASE}?sslmode=disable" -path ./pgsql -verbose up

migrate-down:
	migrate -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_DATABASE}?sslmode=disable" -path ./pgsql -verbose down

migrate-mysql-up:
	migrate -database "mysql://${MYSQL_DB_USERNAME}:${MYSQL_DB_PASSWORD}@tcp(${MYSQL_DB_HOST}:${MYSQL_DB_PORT})/${MYSQL_DB_DATABASE}?charset=utf8&parseTime=true&loc=Local" -path ./mysql -verbose up

migrate-mysql-down:
	migrate -database "mysql://${MYSQL_DB_USERNAME}:${MYSQL_DB_PASSWORD}@tcp(${MYSQL_DB_HOST}:${MYSQL_DB_PORT})/${MYSQL_DB_DATABASE}?charset=utf8&parseTime=true&loc=Local" -path ./mysql -verbose down

migrate-create:
	migrate create -ext sql -seq $(name)