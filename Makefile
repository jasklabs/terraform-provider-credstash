VERSION=v1.1
PACKAGES = $(shell go list ./... | grep -v vendor)

install:
	go install -v

build:
	go build -v -i -o terraform-provider-credstash

test:
	go test $(TESTOPTS) $(PACKAGES)

testacc:
	TF_ACC=1 go test -v $(TESTOPTS) $(PACKAGES) -timeout 120m

release:
	GOOS=darwin go build -v -o darwin_amd64/terraform-provider-credstash_$(VERSION)
	GOOS=linux go build -v -o linux_amd64/terraform-provider-credstash_$(VERSION)
	docker run --rm -v `pwd`:/usr/local/go/src/github.com/sspinc/terraform-provider-credstash -w /usr/local/go/src/github.com/sspinc/terraform-provider-credstash golang:alpine go build -v -o linux_amd64_musl/terraform-provider-credstash_$(VERSION)

.DEFAULT: build
