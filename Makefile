up:
	docker-compose -f docker-compose-prod.yaml up -d --build

down:
	docker-compose down

dev:
	docker-compose up -d --build
