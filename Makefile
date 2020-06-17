PROJECT_NAME := umrs
PKG := github.com/gidyon/$(PROJECT_NAME)
SERVICE_CMD_FOLDER := ${PKG}/cmd
SERVICE_IMAGE_BUILD := ${PKG}/build
SERVICE_PKG_BUILD := ${PKG}/cmd/gateway
API_IN_PATH := api/proto
API_OUT_PATH := pkg/api
SWAGGER_DOC_OUT_PATH := api/swagger
MODULES_DIR := ${PKG}/internal/app
SERVICE_OUT := $(PROJECT_NAME)

ssh: ## Login to aws host
	ssh -i /home/gideon/Desktop/aws/mfcp.pem ubuntu@ec2-54-158-229-126.compute-1.amazonaws.com

setup_dev: ## Sets up a development environment for the umrs project
	@cd deployments/compose/dev &&\
	docker-compose up -d

setup_redis:
	@cd deployments/compose/dev &&\
	docker-compose up -d redis

teardown_dev: ## Tear down development environment for the umrs project
	@cd deployments/compose/dev &&\
	docker-compose down

redis_console:
	@docker run --rm -it --network bridge-rupa-backend redis redis-cli -h redis
	
proto_compile_ledger:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/ledger ledger.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/ledger ledger.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) ledger.proto

proto_compile_patient:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/patient patient.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/patient patient.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) patient.proto

proto_compile_hospital:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/hospital hospital.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/hospital hospital.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) hospital.proto

proto_compile_insurance:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/insurance insurance.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/insurance insurance.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) insurance.proto

proto_compile_account:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/account account.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/account account.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) account.proto

proto_compile_permission:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/permission permission.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/permission permission.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) permission.proto

proto_compile_employment:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/employment employment.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/employment employment.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) employment.proto

proto_compile_treatment:
	protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/treatment treatment.proto &&\
	protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/treatment treatment.proto
	protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) treatment.proto

proto_compile_messaging:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging messaging.proto &&\
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging messaging.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) messaging.proto
	
proto_compile_emailing:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging/emailing emailing.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging/emailing emailing.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) emailing.proto

proto_compile_push:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging/push push.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging/push push.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) push.proto

proto_compile_sms:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging/sms sms.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging/sms sms.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) sms.proto

proto_compile_channel:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging/channel channel.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging/channel channel.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) channel.proto

proto_compile_subscriber:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging/subscriber subscriber.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging/subscriber subscriber.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) subscriber.proto

proto_compile_call:
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/messaging/call call.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/messaging/call call.proto
	protoc -I=$(API_IN_PATH)/messaging -I=third_party --swagger_out=logtostderr=true:$(SWAGGER_DOC_OUT_PATH) call.proto

compile_gateway_prod:
	CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -a -installsuffix cgo -ldflags '-s' -v -o $(SERVICE_OUT) $(SERVICE_PKG_BUILD)

compile_gateway:
	go build -i -v -o gateway $(SERVICE_PKG_BUILD)

TEMPLATE := /home/gideon/go/src/github.com/gidyon/umrs/internal/apps/account/templates
run_account:
	cd  ./cmd/apps/account && go build && export TEMPLATES_DIR=$(TEMPLATE) && ./account -config-file=/home/gideon/go/src/github.com/gidyon/umrs/configs/account.dev.yml

run_notification:
	cd  ./cmd/apps/notification && go build && SMTP_PORT=587 SMTP_PASSWORD=@@antibug2020 SMTP_USERNAME=antibug.ke@gmail.com SMTP_HOST=smtp.gmail.com ./notification -config-file=/home/gideon/go/src/github.com/gidyon/umrs/configs/notification.dev.yml

run_permission:
	REQUEST_ACCESS_TEMPLATE_FILE=/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/permission/templates/request.html PERMISSION_BASE_URL=https://52.14.243.42:30910/ ./cmd/chaincodes/permission/permission -config-file=/home/gideon/go/src/github.com/gidyon/umrs/configs/permission.dev.yml

run_hospital:
	cd ./cmd/chaincodes/hospital/ && go build && TEMPLATES_DIR=/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/hospital/templates ./hospital -config-file=/home/gideon/go/src/github.com/gidyon/umrs/configs/hospital.dev.yml

run_ledger:
	cd ./cmd/ledger && go build && MYSQL_HOST=ec2-18-218-27-110.us-east-2.compute.amazonaws.com MYSQL_USER=root MYSQL_PORT=30760 MYSQL_PASSWORD=@@umrs2020 MYSQL_SCHEMA=umrs-testing ./ledger --sqlite=false

run_ledger_local:
	cd ./cmd/ledger && go build && MYSQL_HOST=localhost MYSQL_USER=root MYSQL_PORT=3306 MYSQL_PASSWORD=hakty11 MYSQL_SCHEMA=umrs ./ledger --sqlite=false

run_insurance:
	ADD_INSURANCE_TEMPLATE_FILE=/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/insurance/templates/add.html DELETE_INSURANCE_TEMPLATE_FILE=/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/insurance/templates/add.html UPDATE_INSURANCE_TEMPLATE_FILE=/home/gideon/go/src/github.com/gidyon/umrs/internal/chaincodes/insurance/templates/add.html  ./cmd/chaincodes/insurance/insurance -config-file=/home/gideon/go/src/github.com/gidyon/umrs/configs/insurance.dev.yml

docker_build_gateway: ## Create a docker image for the service
ifdef tag
	@docker build -t gidyon/$(PROJECT_NAME)-gateway:$(tag) .
else
	@docker build -t gidyon/$(PROJECT_NAME)-gateway:latest .
endif

docker_tag_gateway:
ifdef tag
	@docker tag gidyon/$(PROJECT_NAME)-gateway:$(tag) gidyon/$(PROJECT_NAME)-gateway:$(tag)
else
	@docker tag gidyon/$(PROJECT_NAME)-gateway:latest gidyon/$(PROJECT_NAME)-gateway:latest
endif

docker_push_gateway:
ifdef tag
	@docker push gidyon/$(PROJECT_NAME)-gateway:$(tag)
else
	@docker push gidyon/$(PROJECT_NAME)-gateway:latest
endif

login_cluster:
	ssh -i ~/go/src/github.com/gidyon/umrs/deployments/k8s/keys/cluster/id_rsa admin@api.k8s.coffeeafrica.net

build_and_push_gateway: compile_gateway docker_build_gateway docker_tag_gateway docker_push_gateway

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
