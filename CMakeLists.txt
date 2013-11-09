cmake_minimum_required(VERSION 2.8)

include(FindGit)
include(ExternalProject)

ExternalProject_Add(
	Python3
	URL					http://python.org/ftp/python/3.3.2/Python-3.3.2.tar.bz2
	URL_MD5				7dffe775f3bea68a44f762a3490e5e28
	CONFIGURE_COMMAND	./configure --prefix=${CMAKE_BINARY_DIR}/dist
		# C code messing with sig calls and Go don't mix very well.
		# See https://code.google.com/p/go/issues/detail?id=5287 for details
		COMMAND			cat pyconfig.h | sed s/\#define\ HAVE_SIGALTSTACK\ 1// > pyconfig.new && mv pyconfig.new pyconfig.h
	BUILD_COMMAND		make -j8
	BUILD_IN_SOURCE		1
)

execute_process(
	COMMAND				${GIT_EXECUTABLE} submodule update --init --recursive
	WORKING_DIRECTORY	${CMAKE_SOURCE_DIR}
)

add_custom_command(OUTPUT ${CMAKE_SOURCE_DIR}/3rdparty/libs/gopy/lib/cgo.go
	COMMAND echo "Hello world"
	DEPENDS Python3
)

find_program(GOEXE go)

if(${GOEXE}-NOTFOUND)
	set(GOARCH "386")
	if(CMAKE_SIZEOF_VOID_P EQUAL 8)
		set(GOARCH "amd64")
	endif()
	string(TOLOWER ${CMAKE_SYSTEM_NAME} GOOS)
	set(EXTENSION "tar.gz")
	if(${GOOS} EQUAL "windows")
		set(EXTENSION "zip")
	endif()
	ExternalProject_Add(
		Go
		URL https://go.googlecode.com/files/go1.1.2.${GOOS}-${GOARCH}.${EXTENSION}
		CONFIGURE_COMMAND	""
		BUILD_COMMAND		""
		INSTALL_COMMAND		""
	)

	set(ENV{GOPATH} ${CMAKE_SOURCE_DIR}/../..)
	set(ENV{GOROOT} ${CMAKE_BINARY_DIR}/Go-prefix/src/Go)
	set(GOEXE ${CMAKE_BINARY_DIR}/Go-prefix/src/Go/bin/go)
endif()


foreach(var GOPATH GOROOT GOARCH GOOS)
	execute_process(COMMAND ${GOEXE} env ${var} OUTPUT_VARIABLE ${var} OUTPUT_STRIP_TRAILING_WHITESPACE)
	message("${var}: " ${${var}})
endforeach(var)

set(GOCMD GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOEXE})

add_custom_command(OUTPUT ${GOPATH}/bin/pegparser
	COMMAND ${GOCMD} get github.com/quarnster/parser/pegparser
	DEPENDS ${GOEXE}
)

add_custom_command(OUTPUT ${GOPATH}/pkg/${GOOS}_${GOARCH}/lime/3rdparty/libs/gopy/lib.a
	COMMAND ${GOCMD} install lime/3rdparty/libs/gopy/lib
	DEPENDS ${CMAKE_SOURCE_DIR}/3rdparty/libs/gopy/lib/cgo.go
)

macro(add_peg peg ignore testfile)
	string(REPLACE ".peg" ".go" GOFILE ${peg})
	get_filename_component(dir ${GOFILE} PATH)
	add_custom_command(OUTPUT ${GOFILE}
		COMMAND ${GOPATH}/bin/pegparser -peg=${peg} -outpath=${dir} -notest -ignore=${ignore} -testfile=${testfile}
		DEPENDS ${GOPATH}/bin/pegparser
		DEPENDS ${peg}
		VERBATIM
	)
	set(PEGDEPENDS ${PEGDEPENDS} ${GOFILE})
endmacro(add_peg)

add_peg(${CMAKE_SOURCE_DIR}/backend/loaders/json/json.peg "JsonFile,Values,Value,Null,Dictionary,Array,KeyValuePairs,KeyValuePair,QuotedText,Text,Integer,Float,Boolean,Spacing,Comment" "testdata/Default (OSX).sublime-keymap")
add_peg(${CMAKE_SOURCE_DIR}/backend/loaders/plist/plist.peg "Spacing,KeyValuePair,KeyTag,StringTag,Value,Values,PlistFile,Plist" "testdata/C.plist")

add_custom_command(OUTPUT ${CMAKE_SOURCE_DIR}/backend/sublime/region.go
	COMMAND ${GOCMD} run ${CMAKE_SOURCE_DIR}/build/python.go
	DEPENDS ${GOEXE}
)

add_custom_target(lime
	DEPENDS ${GOPATH}/pkg/${GOOS}_${GOARCH}/lime/3rdparty/libs/gopy/lib.a
	DEPENDS ${CMAKE_SOURCE_DIR}/backend/sublime/region.go
	DEPENDS ${PEGDEPENDS}
)

file(GLOB_RECURSE TESTS ${CMAKE_SOURCE_DIR} "*_test.go")
foreach(testfile ${TESTS})
	if(${testfile} MATCHES ".*._test.go")
		if(NOT ${testfile} MATCHES "bundles|build")
			get_filename_component(dir ${testfile} PATH)
			string(REPLACE ${CMAKE_SOURCE_DIR} "lime" dir ${dir})
			set(GO_TESTS ${GO_TESTS} ${dir})
		endif()
	endif()
endforeach(testfile)

add_custom_target(test
	COMMAND ${GOCMD} test ${GO_TESTS}
	DEPENDS lime
)
