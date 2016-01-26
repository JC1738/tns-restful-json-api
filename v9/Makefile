APP_NAME ?= ot-go-static
PORT0 ?= 8000
LOCALIP=$(shell ifconfig | grep -Eo 'inet (addr:)?([0-9]*\.){3}[0-9]*' | grep -Eo '([0-9]*\.){3}[0-9]*' | grep -v '127.0.0.1' | grep -v '192.168' | grep -v '172.17' | grep -v '172.16' | awk '{print $1}')
TASK_HOST ?= $(LOCALIP)
DISCO_HOST ?= discovery-ci-uswest2.otenv.com
OT_LOGGING_REDIS_HOST ?= logging-qa-uswest2.otenv.com
OT_LOGGING_REDIS_PORT ?= 6379
OT_LOGGING_REDIS_LIST ?= logstash
OT_LOGGING_REDIS_TIMEOUT_MS ?= 5000
STATIC_INSTANCE ?= demo
APP_HOST ?= $(shell hostname)
SHOULD_ANNOUNCE ?= true
DOCKER_RUN_PARAMS ?= --env="PORT0=$(PORT0)" -p="$(PORT0):$(PORT0)" --env="TASK_HOST=$(TASK_HOST)" --env="STATIC_INSTANCE=$(STATIC_INSTANCE)" --dns=10.0.0.104 --dns=10.0.0.103 --dns=10.0.0.102 -v $(PWD)/static:/static 
DOCKERGO = docker run --rm -e CGO_ENABLED=0 -v $(PWD):/usr/src/myapp -v $(PWD):/go -w /usr/src/myapp golang:1.5
buildNumber ?= $(shell git describe --always --tags)
VERSION = $(buildNumber)
REGISTRY ?= docker.otenv.com
APP_VERSION = $(APP_NAME)-$(VERSION)
MESOS_ENV ?= ci-uswest2
FULL_DOCKER_TAG=$(REGISTRY)/$(APP_NAME):$(APP_VERSION)



print-%  : ; @echo $* = $($*)

compile-mac:
	gox -osarch="darwin/amd64" -gcflags="-a" -verbose -output="main_local"

compile-linux:
	gox -osarch="linux/amd64" -gcflags="-a" -verbose -output="main_local"

compile-linux-386:
	gox -osarch="linux/386" -gcflags="-a" -verbose -output="main_local"

local-run:
	TASK_HOST=$(TASK_HOST) APP_HOST=$(APP_HOST) STATIC_INSTANCE=$(STATIC_INSTANCE) ./main_local

compile-run-mac: compile-mac local-run

compile-run-linux: compile-linux local-run

compile-run-linux-386: compile-linux-386 local-run

docker-build:
	docker build --rm -t $(FULL_DOCKER_TAG) .

docker-push:docker-build
	docker tag -f $(FULL_DOCKER_TAG) $(REGISTRY)/$(APP_NAME):latest
	docker push $(REGISTRY)/$(APP_NAME)

docker-run:docker-build
	docker run  $(DOCKER_RUN_PARAMS) $(FULL_DOCKER_TAG)

compile-go:
	rm -rf src
	rm -f main
	cp -r ./Godeps/_workspace/src .
	$(DOCKERGO) go build -a -installsuffix cgo -o main .

compile-run: compile-go docker-run

.PHONY: compile-mac compile-linux compile-linux-386 local-run compile-go compile-run compile-run-mac compile-run-linux compile-run-linux-386 docker-build docker-run docker-push