PROJECT_NAME := umrs
PKG := github.com/gidyon/$(PROJECT_NAME)
SERVICE_CMD_FOLDER := ${PKG}/cmd

compile_binary: ## Compiles ledger binary
	go build -i -v -o ledger $(SERVICE_CMD_FOLDER)/ledger

compile_new: ## Compiles ledger binary with sqlite included
	CGO_ENABLED=1 GOOS=linux go build -o ledger -a -ldflags '-linkmode external -extldflags "-static"' $(SERVICE_CMD_FOLDER)/ledger

compile_binary_prod:
	CGO_ENABLED=1 GOOS=linux GOARCH=386 go build --tags "libsqlite3 linux" -a -installsuffix cgo -ldflags '-s' -v -o ledger $(SERVICE_CMD_FOLDER)/ledger

docker_build: ## Create a docker image for the service
ifdef tag
	@docker build -t gidyon/$(PROJECT_NAME)-ledger:$(tag) .
else
	@docker build -t gidyon/$(PROJECT_NAME)-ledger:latest .
endif

docker_tag: ## Tag a docker image for the service
ifdef tag
	@docker tag gidyon/$(PROJECT_NAME)-ledger:$(tag) gidyon/$(PROJECT_NAME)-ledger:$(tag)
else
	@docker tag gidyon/$(PROJECT_NAME)-ledger:latest gidyon/$(PROJECT_NAME)-ledger:latest
endif

docker_push: ## Push a docker image to a remote repo
ifdef tag
	@docker push gidyon/$(PROJECT_NAME)-ledger:$(tag)
else
	@docker push gidyon/$(PROJECT_NAME)-ledger:latest
endif

compile_and_build_image: compile_new docker_build docker_tag docker_push
build_image: docker_build docker_tag docker_push