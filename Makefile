.PHONY: hep-writer
hep-writer:
	dagger develop
	go build

.PHONY: hep-workspace
hep-workspace:
	dagger --mod .dagger/hep-workspace develop
	go build -C .dagger/hep-workspace

test:
	go test ./...

all: test hep-workspace hep-writer
