DIRS=anomoly exchdb getdata getrss jup
THISDIR=py

all:
	@echo "build completed"

test:
	for dir in $(DIRS); do make -C $$dir test;done
	@echo "test completed in ${THISDIR}"

package:
	for dir in $(DIRS); do make -C $$dir package;done
	@echo "package completed in ${THISDIR}"
