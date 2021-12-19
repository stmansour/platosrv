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

cat > "${MODFILE}" << EOF
CREATE TABLE RSSFeed (
    RSSID BIGINT NOT NULL AUTO_INCREMENT,                   -- unique id for this record
    URL VARCHAR(1024) NOT NULL DEFAULT '',                  -- link to the RSS Feed
    FLAGS BIGINT NOT NULL DEFAULT 0,                        -- no flags defined yet
    LastModTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,  -- when was this record last written
    LastModBy BIGINT NOT NULL DEFAULT 0,                    -- employee UID (from phonebook) that modified it
    CreateTime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,-- when was this record created
    CreateBy BIGINT NOT NULL DEFAULT 0,                     -- employee UID (from phonebook) that created this record
    PRIMARY KEY(RSSID),                                     --
    UNIQUE (URL)                                            -- can't have multiple records with the same Link value
);

CREATE TABLE ItemFeed (
    IFID BIGINT NOT NULL AUTO_INCREMENT,                    -- unique id for this record
    IID BIGINT NOT NULL DEFAULT 0,                          -- The Item
    RSSID BIGINT NOT NULL DEFAULT 0,                        -- The RSSFeed that called it out
    PRIMARY KEY(IFID),
    CONSTRAINT Beta UNIQUE(IID,RSSID)                       -- can't have multiple records with the same IID and RSSID value
);
EOF

#=====================================================
#  Put dir/sqlfilename in the list below
#=====================================================
declare -a dbs=(
	'ws/xb.sql'
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
