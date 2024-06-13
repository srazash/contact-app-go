MAKEFLAGS += --silent

MAIN_PACKAGE_PATH := ./cmd/contact-app
BUILD_PATH := ./build/contact-app

build:
	go build -o ${BUILD_PATH} ${MAIN_PACKAGE_PATH}

run:
	go run ${MAIN_PACKAGE_PATH}

tidy:
	go mod tidy

clean:
	rm -rf ./build/