APP := ws-multicast
OUT := dist

TARGETS := \
	linux/amd64 \
	linux/arm64 \
	linux/arm \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/arm64

.PHONY: all clean build

all: build

build:
	@rm -rf $(OUT)
	@mkdir -p $(OUT)
	@for target in $(TARGETS); do \
		GOOS=$${target%/*}; \
		GOARCH=$${target#*/}; \
		EXT=""; \
		if [ "$$GOOS" = "windows" ]; then EXT=".exe"; fi; \
		echo "Building $$GOOS/$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH \
			go build -trimpath -o "$(OUT)/$(APP)-$$GOOS-$$GOARCH$$EXT" .; \
	done

clean:
	rm -rf $(OUT)

