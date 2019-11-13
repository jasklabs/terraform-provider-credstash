VERSION=v1.1.2
PACKAGES = $(shell go list ./... | grep -v vendor)

install:
	go install -v

build:
	go build -v -i -o build/darwin_amd64/terraform-provider-credstash_$(VERSION)

test:
	go test $(TESTOPTS) $(PACKAGES)

testacc:
	TF_ACC=1 go test -v $(TESTOPTS) $(PACKAGES) -timeout 120m

release:
	GOOS=darwin go build -v -o build/darwin_amd64/terraform-provider-credstash_$(VERSION)
	GOOS=linux go build -v -o build/linux_amd64/terraform-provider-credstash_$(VERSION)
	docker run --rm -v `pwd`:/go/src/github.com/sspinc/terraform-provider-credstash -w /go/src/github.com/sspinc/terraform-provider-credstash golang:alpine go build -v -o build/linux_amd64_musl/terraform-provider-credstash_$(VERSION)

publish:
	aws s3 cp build/darwin_amd64/terraform-provider-credstash_$(VERSION) s3://trident-jars/terraform/darwin_amd64/terraform-provider-credstash_$(VERSION)
	aws s3 cp build/linux_amd64/terraform-provider-credstash_$(VERSION) s3://trident-jars/terraform/linux_amd64/terraform-provider-credstash_$(VERSION)
	aws s3 cp build/linux_amd64_musl/terraform-provider-credstash_$(VERSION) s3://trident-jars/terraform/linux_amd64_musl/terraform-provider-credstash_$(VERSION)

.DEFAULT: build
