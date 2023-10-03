#GOPATH:=$(PWD):${GOPATH}
#export GOPATH
flags=-ldflags="-s -w"
# flags=-ldflags="-s -w -extldflags -static"
TAG := $(shell git tag | sed -e "s,v,,g" | sort -r | head -n 1)

all: build

build:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg; go build -o orecast ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_all: build_darwin_amd64 build_darwin_arm64 build_amd64 build_arm64 build_power8 build_windows

build_darwin_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg orecast_darwin; GOOS=darwin go build -o orecast ${flags}
	mv orecast orecast_darwin_amd64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_darwin_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg orecast_darwin; GOARCH=arm64 GOOS=darwin go build -o orecast ${flags}
	mv orecast orecast_darwin_arm64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_amd64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg orecast_linux; GOOS=linux go build -o orecast ${flags}
	mv orecast orecast_amd64
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

build_power8:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg orecast_power8; GOARCH=ppc64le GOOS=linux go build -o orecast ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv orecast orecast_power8

build_arm64:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg orecast_arm64; GOARCH=arm64 GOOS=linux go build -o orecast ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif
	mv orecast orecast_arm64

build_windows:
ifdef TAG
	sed -i -e "s,{{VERSION}},$(TAG),g" main.go
endif
	go clean; rm -rf pkg orecast.exe; GOARCH=amd64 GOOS=windows go build -o orecast.exe ${flags}
ifdef TAG
	sed -i -e "s,$(TAG),{{VERSION}},g" main.go
endif

install:
	go install

clean:
	go clean; rm -rf pkg

MONGO := $(shell ps auxww | grep mongo | egrep -v grep)

mongo:
ifndef MONGO
	$(error "mongo process not found, please start it to proceed ...")
endif

testdb:
	/bin/rm -f /tmp/files.db && \
	sqlite3 /tmp/files.db < ./schemas/sqlite.sql && \
	mkdir -p /tmp/${USER} && \
	echo "test" > /tmp/${USER}/test.txt

test : mongo testdb test_code

test_code:
	go test -test.v .

# here is an example for execution of individual test
# go test -v -run TestFilesDB
