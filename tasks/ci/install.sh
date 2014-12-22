#!/usr/bin/env bash

# install go get dependencies
sudo apt-get install -qq mercurial

# install backend dependencies
echo 'yes' | sudo add-apt-repository ppa:fkrull/deadsnakes
sudo apt-get update -qq
sudo apt-get install -qq libonig-dev python3.4 python3.4-dev

# install qml frontend dependencies
sudo apt-get install -qq qtbase5-dev qtdeclarative5-dev
