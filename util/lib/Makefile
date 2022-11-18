TOP=../..
THISDIR="util/lib"

lib: config.json *.go
	@touch fail
	go vet
	if [[ -f "${GOPATH}/bin/golint" ]]; then golint; else staticcheck; fi
	go build
	go install
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

config.json:
	@if [ ! -f config.json ]; then \
		/usr/local/bin/jfrog rt dl accord/misc/confdev.json; \
		mv misc/confdev.json . ;\
		rm -rf misc ;\
		cp confdev.json config.json ;\
		fi

clean:
	rm -f log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? ${THISDIR} fail conf*.json
	@echo "*** CLEAN completed in ${THISDIR} ***"

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

test:
	@touch fail
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
