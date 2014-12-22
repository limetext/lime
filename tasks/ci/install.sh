#!/usr/bin/env bash

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
sudo apt-get install -qq qtbase5-dev qtdeclarative5-dev
