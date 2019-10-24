.PHONY: build run psql publish

IMAGE='olblak/polls'
TAG = $(shell git rev-parse HEAD | cut -c1-6)

docker_build:
	docker build -t ${IMAGE}:${TAG} .

docker_publish:
	docker push ${IMAGE}:${TAG}

run:
	bash -c "source sandbox.env && go run main.go"

psql: 
	docker exec -i -t polls_db_1 psql --host=localhost --username=poll -W poll
