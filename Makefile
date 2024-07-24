BINARY_NAME=xte
MAIN_NAME=cmd/main.go
DIST=./dist
FLAGS=-gcflags '-m'

test-coverage:
	go test -v -count=1 -coverprofile=cov.out ./...
	go tool cover -html=cov.out

run:
	go run ${MAIN_NAME} $(arg1)

build:
	GOARCH=amd64 GOOS=darwin go build ${FLAGS} -o ${DIST}-darwin/${BINARY_NAME} ${MAIN_NAME}
	GOARCH=amd64 GOOS=linux go build ${FLAGS} -o ${DIST}-linux/${BINARY_NAME} ${MAIN_NAME}
	GOARCH=amd64 GOOS=windows go build ${FLAGS} -o ${DIST}-windows/${BINARY_NAME}.exe ${MAIN_NAME}

deps:
	@[ -d $(DIST)-linux ] || mkdir -p $(DIST)-linux
	@[ -d $(DIST)-darwin ] || mkdir -p $(DIST)-darwin
	@[ -d $(DIST)-windows ] || mkdir -p $(DIST)-windows

clean:
	@rm -rf ${DIST}*
	@-rm -f *.pprof *.prof *.test *.out

all: deps build

# utility for printing variables
print-% : ; @echo $($*)