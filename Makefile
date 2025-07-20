default: all
all: test hep-workspace hepgen

.PHONY: hepgen
hepgen:
	dagger develop
	go build

.PHONY: hep-workspace
hep-workspace:
	dagger --mod .dagger/hep-workspace develop
	go build -C .dagger/hep-workspace

test:
	go test ./...

hep:
	dagger -c /bin/sh -c 'hep'

preview:
	dagger -c /bin/sh -c 'preview|up'

sandbox:
	dagger -c /bin/sh -c 'sandbox|up'
