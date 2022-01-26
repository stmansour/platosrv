#!/usr/bin/env bash
#
#=============================================================================


#-------------------------------------------------------------------------------
#  usage  -  Display information explaining how to use this script
#-------------------------------------------------------------------------------
usage() {
    PROGNAME="collector.sh"
    cat <<ZZEOF

Plato DB Data Collector

    Usage:   ${PROGNAME} [OPTIONS]

    This command updates the mysql database named "plato" with exchange
    and RSS feed data. If no option is provided, it will update the database
    with today's exchange and RSS feed information. Or, with the -update option,
    it will update all the information in the database since the last update.

OPTIONS
    -u, -update  will update the plato database with all the RSS feed and
    exchange data since the last detected update.


EXAMPLES:
    Command to update plato database with today's RSS feeds and yesterday's
    currency exchange rates:

        bash$  ./${PROGNAME}

    Command to update plato database with all information since the last time
    this script was run:

        bash$  ./${PROGNAME} -update

ZZEOF
}

Today() {
    ./edb.sh
    ./collect.sh
}

Update() {
    ./edb.sh update
    ./wbp.sh update
}

for arg do
	# echo '--> '"\`$arg'"
	cmd=$(echo "${arg}" |tr "[:upper:]" "[:lower:]")
    case "$cmd" in
    "debug" | "-debug" | "-d" )
        DEBUG=1
        ;;
    "help" | "h" | "-h" | "-help")
        usage
        exit 0
        ;;
    "-u" | "-update" | "u" | "update")
        Update
        ;;

	*)  #invalid argument
		echo "Unrecognized command: $arg"
		usage
		exit 1
		;;
    esac
done

Today
