THISDIR=py/getexch
TOPDIR=../..
DIST=${TOPDIR}/dist/plato

getdata:
	echo "**** Finished in ${THISDIR}"

clean:
	rm -rf xxx tmprss completed.txt doit.sh

cleanall: clean
	rm *.csv

package:
	mkdir -p ${DIST}
	cp *.py *.sh *.txt ${DIST}/

test:
	echo "Done in ${THISDIR}"
