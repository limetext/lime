#!/bin/sh
set -e

if [ "$TRAVIS_OS_NAME" = "linux" ]; then
  echo 'yes' | sudo add-apt-repository ppa:fkrull/deadsnakes
  sudo apt-get update -qq
  sudo apt-get install -qq libonig-dev mercurial python3.4 python3.4-dev xclip
elif [ "$TRAVIS_OS_NAME" = "osx" ]; then
  sudo brew install pkg-config go mercurial oniguruma python3
  export PKG_CONFIG_PATH=$(brew --prefix python3)/Frameworks/Python.framework/Versions/3.4/lib/pkgconfig
else
  echo "BUILD NOT CONFIGURED: $TRAVIS_OS_NAME"
  exit 1
fi
