#!/usr/bin/env bash

if [ "$TRAVIS_OS_NAME" = "linux" ]; then

	# make sure we're up to date
	sudo apt-get update -qq

	# install go get dependencies
	sudo apt-get install -qq mercurial

	# install backend dependencies
	sudo add-apt-repository -y ppa:fkrull/deadsnakes
	sudo apt-get update -qq
	sudo apt-get install -qq libonig-dev python3.4 python3.4-dev

	# install qml frontend dependencies
	sudo add-apt-repository -y ppa:ubuntu-sdk-team/ppa
	sudo apt-get update -qq
	sudo apt-get install -qq qtbase5-private-dev qtdeclarative5-private-dev

elif [ "$TRAVIS_OS_NAME" = "osx" ]; then

	brew install pkg-config go mercurial oniguruma python3
	export PKG_CONFIG_PATH=$(brew --prefix python3)/Frameworks/Python.framework/Versions/3.4/lib/pkgconfig

else

	echo "BUILD NOT CONFIGURED: $TRAVIS_OS_NAME"
	exit 1

fi
