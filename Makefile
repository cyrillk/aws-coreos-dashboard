default: build

NAME=aws-coreos-dashboard

build:
	go build .	

docker:
	docker build -t $(NAME) .
