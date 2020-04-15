# Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved. 
# SPDX-License-Identifier: Apache-2.0

.PHONY: build clean test run docker update

GO=CGO_ENABLED=0 GO111MODULE=on GO
CGO=CGO_ENABLED=1 GO111MODULE=on GO

MICROSERVICES=cmd/edgex-club/edgex-club
DOCKER=docker-edgex-club

.PHONY: $(MICROSERVICES) $(DOCKER)


VERSION=$(shell cat ./VERSION 2>/dev/null || echo 0.0.0)

GOFLAGS=-ldflags "-X edgex-club.Version=$(VERSION)"
GIT_SHA=$(shell git rev-parse HEAD)

build: $(MICROSERVICES)
	$(GO) build ./...

cmd/edgex-club/edgex-club: 
	$(GO) build $(GOFLAGS) -o $@ ./cmd/edgex-club

clean:
	rm -f $(MICROSERVICES)

test:
	GO111MODULE=on go test -coverprofile=coverage.out ./...
	GO111MODULE=on go vet ./...

update:
	$(GO) mod download

run:
	cd bin && ./edgex-club-launch.sh

docker: $(DOCKERS)

docker-edgex-club: 
	docker build --label "git_sha=$(GIT_SHA)" -t edgexclub/docker-edgex-club:$(VERSION) .