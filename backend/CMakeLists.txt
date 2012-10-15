cmake_minimum_required(VERSION 2.6)

include_directories(3rdparty/PlistCpp/src)
add_library(backend STATIC
    3rdparty/PlistCpp/src/pugixml.cpp
    Syntax.cpp
)
