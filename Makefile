docker:
	docker-compose -f docker-compose.yml up -d

.DEFAULT_GOAL := docker


stop: 
	docker-compose stop

db:
	docker exec -it auth-psql psql auth_user -d auth_db