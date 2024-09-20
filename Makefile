BINARY_NAME=xte
MAIN_NAME=cmd/main.go
DIST=./dist
FLAGS=-gcflags '-m'

test-perf:
	go test -count=1  -c ./util && util.test -test.count=1 -test.benchmem  -test.bench=. -test.cpuprofile cpu.prof -test.benchtime 1000s

test-coverage:
	go test -v -count=1 -coverprofile=cov.out ./...
	go tool cover -html=cov.out

run:
	go run ${MAIN_NAME} $(arg1)

build-windows: clean
	@[ -d $(DIST)-windows ] || mkdir -p $(DIST)-windows
	GOARCH=amd64 GOOS=windows go build ${FLAGS} -o ${DIST}-windows/${BINARY_NAME}.exe ${MAIN_NAME}

build: deps
	go build ${FLAGS} -o ${DIST}/${BINARY_NAME} ${MAIN_NAME}

build-static: deps
	CGO_ENABLED=0 go build -o ${DIST}/${BINARY_NAME} -a -ldflags="-w -extldflags" ${MAIN_NAME}

deps:
	@[ -d $(DIST) ] || mkdir -p $(DIST)

clean:
	@rm -rf ${DIST}*
	@-rm -f *.pprof *.prof *.test *.out

all: clean deps build

# utility for printing variables
print-% : ; @echo $($*)