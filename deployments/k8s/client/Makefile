PROJECT_NAME := umrs
PKG := github.com/gidyon/$(PROJECT_NAME)
SERVICE_CMD_FOLDER := ${PKG}/cmd/ledger/client

compile_binary: ## Compiles ledger-client binary
	go build -i -v -o ledger-client $(SERVICE_CMD_FOLDER)

docker_build: ## Create a docker image for the service
ifdef tag
	@docker build -t gidyon/$(PROJECT_NAME)-ledger-client:$(tag) .
else
	@docker build -t gidyon/$(PROJECT_NAME)-ledger-client:latest .
endif

docker_tag: ## Tag a docker image for the service
ifdef tag
	@docker tag gidyon/$(PROJECT_NAME)-ledger-client:$(tag) gidyon/$(PROJECT_NAME)-ledger-client:$(tag)
else
	@docker tag gidyon/$(PROJECT_NAME)-ledger-client:latest gidyon/$(PROJECT_NAME)-ledger-client:latest
endif

docker_push: ## Push a docker image to a remote repo
ifdef tag
	@docker push gidyon/$(PROJECT_NAME)-ledger-client:$(tag)
else
	@docker push gidyon/$(PROJECT_NAME)-ledger-client:latest
endif

bundle: compile_binary docker_build docker_tag docker_push