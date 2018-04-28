
build:
	docker build . -t reddit-watch:latest

run:
	docker run --rm -it reddit-watch:latest
