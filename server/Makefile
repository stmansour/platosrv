TOP=..
BINDIR=${TOP}/dist/plato
THISDIR="server"
COUNTOL=${TOP}/tools/bashtools/countol.sh

server: *.go config.json
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	@rm -f fail
	@echo "*** ${THISDIR}: completed ${THISDIR} ***"

# you can update the config files as follows:
# $ jfrog rt u[pload] configplatosrvdev.json accord/misc/
# $ deployfile.sh configplatosrvdev.json accord/db/

config.json:
	@if [ -f ${TOP}/util/lib/config.json ]; then \
		cp ${TOP}/util/lib/config.json . ;\
	else \
		/usr/local/bin/jfrog rt dl accord/misc/confdev.json ;\
		mv misc/confdev.json . ;\
		rm -rf misc ;\
		cp confdev.json config.json ;\
	fi

pushconf:
	jfrog rt u confdev.json accord/misc/
	/usr/local/accord/bin/deployfile.sh confdev.json accord/db/

package: server
	mkdir -p ${BINDIR}
	cp server ${BINDIR}/platosrv
	cp confdev.json ${BINDIR}/config.json

build: server package

ver.go:
	./mkver.sh

clean:
	go clean
	rm -f ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean server package

secure:
	@rm -f config.json configplatosrvdev.json confprod.json
