
build:
	docker build . -t reddit-watch:latest

run:
	docker run --rm -it -v `pwd`/config:/config reddit-watch:latest --token ${TOKEN}
