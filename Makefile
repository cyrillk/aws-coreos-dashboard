default: build

NAME=aws-coreos-dashboard
REGISTRY=cyrillk/

build:
	docker build --force-rm -t $(NAME) .

push:
	docker tag $(NAME) $(REGISTRY)$(NAME)
	docker push $(REGISTRY)$(NAME)
