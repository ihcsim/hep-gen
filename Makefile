default: all
all: test hep-workspace hep-gen

.PHONY: hep-gen
hep-gen:
	dagger develop
	go build

.PHONY: hep-workspace
hep-workspace:
	dagger --mod .dagger/hep-workspace develop
	go build -C .dagger/hep-workspace

test:
	go test ./...

hep:
	dagger -c /bin/sh -c 'hep "test title" | terminal'

preview:
	dagger -c /bin/sh -c 'preview| as-service --args "madness" --args "server" | up'

sandbox:
	dagger -c /bin/sh -c 'sandbox| terminal'
