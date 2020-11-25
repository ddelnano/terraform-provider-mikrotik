.PHONY: import testacc dist

TEST ?= ./...

build:
	go build -o terraform-provider-mikrotik

clean:
	rm dist/*
	terraform-provider-mikrotik

plan: build
	terraform init
	terraform plan

apply:
	terraform apply

test:
	go test $(TEST) -v

testacc:
	TF_ACC=1 go test $(TEST) -count 1 -v
