ifndef ${GOPATH}
GOPATH=$(shell pwd)
export GOPATH
endif

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
SRC=src/
MODELS=webserver/models
ROUTERS=webserver
DATABASE=webserver/dbconnection
CONTROLLERS=webserver/controllers

$(shell cd $(SRC))

all: model controller database router
	${GOBUILD} start.go

model:
	${GOBUILD} ${MODELS}

controller:
	${GOBUILD}  ${CONTROLLERS}

database:
	${GOBUILD}  ${DATABASE}

router:
	${GOBUILD} ${ROUTERS}

.PHONY: test open install

test:
	${GOTEST} ${MODELS}
	${GOTEST} ${CONTROLLERS}
	# ${GOTEST} ${SRC}

open:
	$(shell sudo setcap cap_net_bind_service=+ep `pwd`/start)

install:
	${GOINSTALL} ${MODELS}
	${GOINSTALL} ${CONTROLLERS}
	${GOINSTALL} ${DATABASE}
	${GOINSTALL} ${ROUTERS}
