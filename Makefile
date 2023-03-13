build:
	docker-compose up --build -d

up:
	docker-compose up -d

bash:
	docker-compose exec airflow-webserver bash

restart:
	docker-compose restart

exec:
	docker-compose exec click_server clickhouse-client

down:
	docker-compose down

kafka:
	docker exec -it kafka_fungjai bash