#!/usr/bin/env bash

while [ $# != 0 ]; do
	arg="$1"
	case "$arg" in
		-f|--force)
			force=true
			;;
	esac
	shift
done

if [ "$TRAVIS" != "true" ] && [ "$force" != "true" ]; then
	echo "WARNING: This script will destroy any unstaged changes."
	echo "Do you want to continue? [Y/n]"
	read cont
	if [ "$cont" != "" ] && [ "$cont" != "y" ] && [ "$cont"  != "Y" ]; then
		exit 1
	fi
fi
