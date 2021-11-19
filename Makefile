DIRS=util session db
DIST=dist
.PHONY: test

plato:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	rm -rf dist

test:
	for dir in $(DIRS); do make -C $$dir test;done

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean plato package test stats
	echo "Completed"

build: clean plato package

release:
	/usr/local/accord/bin/release.sh plato

tarzip:
	cd ${DIST};if [ -f ./plato/config.json ]; then mv ./plato/config.json .; fi
	cd ${DIST};rm -f plato.tar*;tar czf plato.tar.gz plato
	cd ${DIST};if [ -f ./config.json ]; then mv ./config.json ./plato/config.json; fi

stats:
	@echo
	@echo "-------------------------------------------------------------------------------"
	@echo "|                         GO SOURCE CODE STATISTICS                           |"
	@echo "-------------------------------------------------------------------------------"
	@find . -name "*.go" | srcstats
	@echo "-------------------------------------------------------------------------------"
