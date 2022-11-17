TOP=../..
BINDIR=${TOP}/dist/plato
THISDIR="admin/dbcheck"
COUNTOL=${TOP}/tools/bashtools/countol.sh

dbcheck: config.json *.go
	@find . -name "fail" -exec rm -r "{}" \;
	@touch fail
	golint
	@if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	# go install
	@rm -f fail
	@echo "*** ${THISDIR}: completed ws ***"

config.json:
	@if [ -f ${TOP}/util/lib/config.json ]; then \
		cp ${TOP}/util/lib/config.json . ;\
	else \
		jf rt dl accord/misc/confdev.json ;\
		mv misc/confdev.json . ;\
		rm -rf misc ;\
		cp confdev.json config.json ;\
	fi

clean:
	go clean
	rm -f ver.go tmp *.log *.json
	@echo "*** ${THISDIR}: completed clean ***"

test:
	@echo "*** ${THISDIR}: completed test ***"

all: clean dbcheck package

mkprod:
	/usr/local/accord/bin/getfile.sh accord/db/confprod.json
	cp confprod.json config.json
secure:
	@rm -f config.json confdev.json confprod.json

package:
	cp dbcheck ${BINDIR}/
	@echo "*** ${THISDIR}: completed package ***"
