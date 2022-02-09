TOP=../..
THISDIR="dbtest"

dbtest: *.go config.json
	@touch fail
	@go vet
	@if [[ `uname` = "Darwin" ]]; then staticcheck; else golint; fi
	go build
	@rm -f fail
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f conf*.json log llog *.g ./gold/*.g err.txt [a-z] [a-z][a-z1-9] qq? ${THISDIR} fail
	@echo "*** CLEAN completed in ${THISDIR} ***"

config.json:
	jfrog rt dl accord/misc/confdev.json
	mv misc/confdev.json .; cp confdev.json config.json
	rm -rf misc

relink:
	go build
	@echo "*** RELINK completed in ${THISDIR} ***"

test: dbtest
	@touch fail
	mysql --no-defaults plato  < ../schema/db.sql
	./dbtest
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
