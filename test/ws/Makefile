THISDIR="ws"
TOP=../..
all:
	@echo "*** Completed in ${THISDIR} ***"

clean:
	rm -f loginrequest platosrv.log slog wreis.log serverreply log llog *.g ./gold/*.g err.txt [a-z] [a-z][0-9] [a-z0-9][0-9][0-9] qq? fail request w.[a-z] con*.json
	@echo "*** CLEAN completed in ${THISDIR} ***"

test: config.json
	@touch fail
	./functest.sh
	@echo "*** TEST completed in ${THISDIR} ***"
	@rm -f fail

config.json:
	@if [ -f ${TOP}/util/lib/config.json ]; then \
		cp ${TOP}/util/lib/config.json . ;\
	else \
		/usr/local/bin/jfrog rt dl accord/misc/confdev.json ;\
		mv misc/confdev.json . ;\
		rm -rf misc ;\
		cp confdev.json config.json ;\
	fi

package:
	@echo "*** PACKAGE completed in ${THISDIR} ***"

secure:
	@rm -f config.json confdev.json confprod.json
