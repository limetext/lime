cmake_minimum_required(VERSION 2.8)

include(FindGit)
include(ExternalProject)

if(WIN32)
	set(EXESUFFIX ".exe")
	find_package(PythonInterp 3.3 QUIET)
	if (NOT PYTHONINTERP_FOUND)
		set(file python-3.3.2.msi)
		if(CMAKE_SIZEOF_VOID_P EQUAL 8)
			set(file python-3.3.2.amd64.msi)
		endif()
		if(NOT EXISTS ${CMAKE_BINARY_DIR}/${file})
			message("Downloading python ${file}")
			file(DOWNLOAD http://python.org/ftp/python/3.3.2/${file} ${CMAKE_BINARY_DIR}/${file} SHOW_PROGRESS EXPECTED_MD5 0d9db9c2316562c62e1e4c347b6f9430)
		endif()
		set(PYTHON_PATH "C:\\Python3")
		if(NOT EXISTS ${PYTHON_PATH})
			find_program(MSIEXEC msiexec REQUIRED)
			message("Installing python3 to ${PYTHON_PATH}")
			execute_process(
				COMMAND ${MSIEXEC} /i ${file} TARGETDIR=${PYTHON_PATH} /qb /log test.log
				WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
			)
		endif()
		add_custom_command(OUTPUT ${CMAKE_SOURCE_DIR}/3rdparty/libs/gopy/lib/cgo.go
			COMMAND echo "Hello world"
		)
	endif()
else()
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
	add_custom_command(OUTPUT ${CMAKE_SOURCE_DIR}/3rdparty/libs/gopy/lib/cgo.go
		COMMAND echo "Hello world"
		DEPENDS Python3
	)
endif()

# ExternalProject_Add(
# 	Oniguruma
# 	URL			http://www.geocities.jp/kosako3/oniguruma/archive/onig-5.9.5.tar.gz
# 	URL_MD5		970f98a4cd10021b545d84e34c34aae4
# )

execute_process(
	COMMAND				${GIT_EXECUTABLE} submodule update --init --recursive
	WORKING_DIRECTORY	${CMAKE_SOURCE_DIR}
)


find_program(GOEXE go)
if(NOT GOEXE)
	set(GOARCH "386")
	if(CMAKE_SIZEOF_VOID_P EQUAL 8)
		set(GOARCH "amd64")
	endif()
	string(TOLOWER ${CMAKE_SYSTEM_NAME} GOOS)
	set(EXTENSION "tar.gz")
	if(GOOS STREQUAL "windows")
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
	set(GOEXE ${CMAKE_BINARY_DIR}/Go-prefix/src/Go/bin/go${EXESUFFIX})
endif()


foreach(var GOPATH GOROOT GOARCH GOOS)
	execute_process(COMMAND ${GOEXE} env ${var} OUTPUT_VARIABLE ${var} OUTPUT_STRIP_TRAILING_WHITESPACE)
	message("${var}: " ${${var}})
endforeach(var)


if(WIN32)
	file(WRITE "${CMAKE_BINARY_DIR}/go.bat" "
set GOPATH=${GOPATH}
set GOROOT=${GOROOT}
${GOEXE} %*
")
	set(GOCMD ${CMAKE_BINARY_DIR}/go.bat)
else()
	set(GOCMD GOROOT=${GOROOT} GOPATH=${GOPATH} ${GOEXE})
endif()

set(GOPKG ${GOPATH}/pkg/${GOOS}_${GOARCH})
foreach(pkg code.google.com/p/log4go github.com/quarnster/parser github.com/quarnster/util/text)
	add_custom_command(OUTPUT ${GOPKG}/${pkg}.a
		COMMAND ${GOCMD} get ${pkg}
		DEPENDS ${GOEXE}
	)
endforeach()

add_custom_command(OUTPUT ${GOPATH}/bin/pegparser${EXESUFFIX}
	COMMAND ${GOCMD} install github.com/quarnster/parser/pegparser
	DEPENDS ${GOEXE}
	DEPENDS ${GOPKG}/github.com/quarnster/parser.a
)

add_custom_command(OUTPUT ${GOPKG}/lime/3rdparty/libs/gopy/lib.a
	COMMAND ${GOCMD} install lime/3rdparty/libs/gopy/lib
	DEPENDS ${CMAKE_SOURCE_DIR}/3rdparty/libs/gopy/lib/cgo.go
)

macro(add_peg peg ignore testfile)
	string(REPLACE ".peg" ".go" GOFILE ${peg})
	get_filename_component(dir ${GOFILE} PATH)
	add_custom_command(OUTPUT ${GOFILE}
		COMMAND ${GOPATH}/bin/pegparser -peg=${peg} -outpath=${dir} -notest -ignore=${ignore} -testfile=${testfile}
		DEPENDS ${GOPATH}/bin/pegparser${EXESUFFIX}
		DEPENDS ${peg}
		WORKING_DIRECTORY ${CMAKE_BINARY_DIR}
		VERBATIM
	)
	set(PEGDEPENDS ${PEGDEPENDS} ${GOFILE})
endmacro(add_peg)

add_peg(${CMAKE_SOURCE_DIR}/backend/loaders/json/json.peg "JsonFile,Values,Value,Null,Dictionary,Array,KeyValuePairs,KeyValuePair,QuotedText,Text,Integer,Float,Boolean,Spacing,Comment" "testdata/Default (OSX).sublime-keymap")
add_peg(${CMAKE_SOURCE_DIR}/backend/loaders/plist/plist.peg "Spacing,KeyValuePair,KeyTag,StringTag,Value,Values,PlistFile,Plist" "testdata/C.plist")

add_custom_command(OUTPUT ${CMAKE_SOURCE_DIR}/backend/sublime/region.go
	COMMAND ${GOCMD} run ${CMAKE_SOURCE_DIR}/build/python.go
	DEPENDS ${PEGDEPENDS}
)

add_custom_target(lime
	DEPENDS ${CMAKE_SOURCE_DIR}/backend/sublime/region.go
	DEPENDS ${GOPATH}/pkg/${GOOS}_${GOARCH}/lime/3rdparty/libs/gopy/lib.a
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
