VERSION ?= vlatest

TARGETS := amd64 arm6 arm7

ARCH_arm6 := GOARCH=arm GOARM=6
ARCH_arm7 := GOARCH=arm GOARM=7
ARCH_amd64 := GOARCH=amd64

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
PROJECT_DIR=${CURDIR}
CMDS = $(shell ls ${PROJECT_DIR}/cmd)
DIST_DIR=${PROJECT_DIR}/dist/$@

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

all: $(TARGETS)

$(TARGETS): clean test
	mkdir ${DIST_DIR};
	mkdir ${DIST_DIR}/cmd;
	mkdir ${DIST_DIR}/storage;
	mkdir -p ${DIST_DIR}/assets/config;
	mkdir ${DIST_DIR}/bin;
	cp ${PROJECT_DIR}/resources -r ${DIST_DIR};
	cp ${PROJECT_DIR}/.env.example ${DIST_DIR}/.env;

	${ARCH_$@} go build ${LDFLAGS} -o ${DIST_DIR}/goairmon;
	${ARCH_$@} go build ${LDFLAGS} -o ${DIST_DIR}/cmd/adduser ./cmd/adduser;

	#cp scripts/install.sh scripts/uninstall.sh scripts/lasecap.service ${DIST_DIR}/bin
	#cp scripts/start_hotspot.sh scripts/stop_hotspot.sh ${DIST_DIR}/bin

	tar -czf ${DIST_DIR}/../goairmon-$@.tar.gz -C ${DIST_DIR}/ .

test:
	cd ${PROJECT_DIR}
	go test ./...

clean:
	cd ${PROJECT_DIR};
	rm -rf dist
	mkdir dist