rebuild:
	docker-compose up --build $(word 2,$(MAKECMDGOALS)) -d
	docker-compose logs -f $(word 2,$(MAKECMDGOALS))

rebuild-all:
	docker-compose up --build -d

log:
	docker-compose logs -f $(word 2,$(MAKECMDGOALS))

restart:
	docker-compose restart $(word 2,$(MAKECMDGOALS))

build:
	docker-compose build $(word 2,$(MAKECMDGOALS))

up:
	docker-compose up -d

down:
	docker-compose down

clean:
	docker-compose down -v

%:
	@:
