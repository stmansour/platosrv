#!/bin/bash

#==========================================================================
#  This script performs SQL schema changes on the test databases that are
#  saved as SQL files in the test directory. It loads them, performs the
#  ALTER commands, then saves the sql file.
#
#  If the test file uses its own database saved as a .sql file, make sure
#  it is listed in the dbs array
#==========================================================================

MODFILE="dbqqqmods.sql"
MYSQL="mysql --no-defaults"
MYSQLDUMP="mysqldump --no-defaults"
DBNAME="plato"

#=====================================================
#  Retain prior changes as comments below
#=====================================================

#=====================================================
#  Put modifications to schema in the lines below
#=====================================================

cat > "${MODFILE}" << LEOF
LEOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	# 'ws/xb.sql'
	# 'ws/xh.sql'
)

for f in "${dbs[@]}"
do
	echo "DROP DATABASE IF EXISTS plato; CREATE DATABASE plato; USE plato; GRANT ALL PRIVILEGES ON plato.* TO 'ec2-user'@'localhost';" | ${MYSQL}
	echo -n "${f}: loading... "
	${MYSQL} ${DBNAME} < ${f}
	echo -n "updating... "
	${MYSQL} ${DBNAME} < ${MODFILE}
	echo -n "saving... "
	${MYSQLDUMP} ${DBNAME} > ${f}
	echo "done"
done
