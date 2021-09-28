.PHONY: build
build: build/builder build/source

.PHONY: build/builder
build/builder:
	docker build -t gtk3-builder:latest docker/

.PHONY: run
run: build/builder
	docker run \
		--rm \
		--env DISPLAY \
		--volume gopkgcache:/root/go/pkg \
		--volume gobuildcache:/root/.cache/go-build \
		--volume cinnycache:/root/.cache/cinny-desktop \
		--mount type=bind,source=/tmp/.X11-unix,target=/tmp/.X11-unix \
		--mount type=bind,source=/var/run/dbus/system_bus_socket,target=/var/run/dbus/system_bus_socket \
		--mount type=bind,source=$(shell pwd),target=/source \
		gtk3-builder:latest \
		make -C cmd/cinny-desktop run

# TODO: pin specific cinny version
.PHONY: download/cinny
download/cinny: pkg/assets/cinny/index.html

pkg/assets/cinny/index.html:
	bash download_cinny.sh
