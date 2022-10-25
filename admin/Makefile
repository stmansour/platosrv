DIRS=dbcheck

main:
	for dir in $(DIRS); do make -C $$dir;done

clean:
	for dir in $(DIRS); do make -C $$dir clean;done

package:
	for dir in $(DIRS); do make -C $$dir package;done

all: clean main package
	@echo "Completed"

build: clean main package

