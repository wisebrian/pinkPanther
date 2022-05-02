GIT_SHA=$(shell git rev-parse --short HEAD)

go-install:
	CGO_ENABLED=0 go install ./...

go-build:
	CGO_ENABLED=0 go build ./...

go-test:
	go test ./...

build:
	docker build -f Dockerfile -t revox:$(GIT_SHA) -t revox:latest .

start:
	docker-compose up -d

restart:
	docker-compose build --no-cache
	docker-compose up -d --force-recreate

stop:
	docker-compose stop

state:
	docker-compose ps

logs:
	docker-compose logs

# Helm stuff

install:
	helm upgrade -i revox ./helm

template:
	helm template revox ./helm --values ./helm/values-demo.yaml