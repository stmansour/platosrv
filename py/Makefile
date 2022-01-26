DIRS=getexch getrss
THISDIR=py

all:
	@echo "build completed"

clean:
	for dir in $(DIRS); do make -C $$dir clean;done
	@echo "clean completed in ${THISDIR}"

test:
	for dir in $(DIRS); do make -C $$dir test;done
	@echo "test completed in ${THISDIR}"

package:
	for dir in $(DIRS); do make -C $$dir package;done
	@echo "package completed in ${THISDIR}"
