cmake_minimum_required(VERSION 2.6)

set(CURSES_NEED_NCURSES True)
find_package(curses REQUIRED) # Required for now..

if(CURSES_FOUND)
    set(CURSES_FRONTEND frontend/curses/main.cpp)
    include_directories(${CURSES_INCLUDE_DIR})
endif()
set(FRONTENDS ${CURSES_FRONTEND})


include(FindPackageHandleStandardArgs)
macro(my_find_package NAME HEADER)
    string(TOUPPER ${NAME} NU)
    find_library(${NU}_LIBRARY ${NAME})
    find_path(${NU}_INCLUDE_DIR ${HEADER})
    find_package_handle_standard_args(${NAME} REQUIRED_VARS ${NU}_INCLUDE_DIR ${NU}_LIBRARY)
    mark_as_advanced(${NU}_INCLUDE_DIR ${NU}_LIBRARY)
endmacro()


my_find_package(Carbon Carbon/Carbon.h)
if(NOT CARBON_FOUND)
    message(FATAL_ERROR "Carbon not found")
endif()
my_find_package(termkey termkey.h)
if(TERMKEY_FOUND)
    add_definitions(-DHAVE_TERMKEY)
endif()

include_directories(backend)

add_executable(lime ${FRONTENDS} ${BACKEND})
target_link_libraries(lime ${CURSES_LIBRARIES} ${CARBON_LIBRARY} ${TERMKEY_LIBRARY})
