.PHONY: build generate clean plan apply lint-client lint-provider lint testacc testclient test

TIMEOUT ?= 40m
ROUTEROS_VERSION ?= "6.48.3"
ifdef TEST
    override TEST := ./... -run $(TEST)
else
    override TEST := ./...
endif

ifdef TF_LOG
    override TF_LOG := TF_LOG=$(TF_LOG)
endif

build:
	go build -o terraform-provider-mikrotik

generate:
	go generate ./...

clean:
	rm dist/*

plan: build
	terraform init
	terraform plan

apply:
	terraform apply

lint-client:
	go vet ./client/...

lint-provider:
	go vet ./mikrotik/...

lint: lint-client lint-provider

test: lint testclient testacc

testclient:
	cd client; go test $(TEST) -race -v -count 1

testacc:
	TF_ACC=1 $(TF_LOG) go test $(TEST) -v -count 1 -timeout $(TIMEOUT)

routeros: routeros-clean
	ROUTEROS_VERSION=$(ROUTEROS_VERSION) docker compose -f docker/docker-compose.yml up -d routeros

routeros-stop:
	docker compose -f docker/docker-compose.yml stop routeros

routeros-clean:
	docker compose -f docker/docker-compose.yml rm -sfv routeros
