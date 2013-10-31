#!/usr/bin/env sh

wget http://python.org/ftp/python/3.3.2/Python-3.3.2.tar.bz2
tar -xjf Python-3.3.2.tar.bz2
cd Python-3.3.2
./configure --enable-shared
cat pyconfig.h | sed s/#define\ HAVE_SIGALTSTACK\ 1// > pyconfig.new && mv pyconfig.new pyconfig.h
make -j8
sudo make install
