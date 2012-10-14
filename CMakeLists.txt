cmake_minimum_required(VERSION 2.6)

set(CURSES_NEED_NCURSES True)
find_package(curses REQUIRED) # Required for now..

if(CURSES_FOUND)
    set(CURSES_FRONTEND frontend/curses/main.cpp)
    include_directories(${CURSES_INCLUDE_DIR})
endif()
set(FRONTENDS ${CURSES_FRONTEND})

find_library(APP_SERVICES_LIBRARY ApplicationServices)
find_library(CORE_SERVICES_LIBRARY CoreServices)
find_library(CARBON_LIBRARY Carbon)
mark_as_advanced(APP_SERVICES_LIBRARY CORE_SERVICES_LIBRARY CARBON_LIBRARY)
if(NOT CARBON_LIBRARY)
    message(FATAL_ERROR "Carbon not found")
endif()

add_executable(lime ${FRONTENDS} ${BACKEND})
target_link_libraries(lime ${CURSES_LIBRARIES} ${APP_SERVICES_LIBRARY} ${CORE_SERVICES_LIBRARY} ${CARBON_LIBRARY})
