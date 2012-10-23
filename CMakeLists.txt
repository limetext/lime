cmake_minimum_required(VERSION 2.6)

find_package(Qt4 COMPONENTS QtGui QtCore REQUIRED)
find_package(Boost COMPONENTS Python REQUIRED)
find_package(PythonLibs REQUIRED)
include(${QT_USE_FILE})
add_definitions(${QT_DEFINITIONS})
include_directories(${PYTHON_INCLUDE_DIRS} ${Boost_INCLUDE_DIRS})

qt4_wrap_cpp(lime_MOC frontend/qt/MainLoop.h)
add_executable(lime MACOSX_BUNDLE
    frontend/qt/main.cpp
    frontend/qt/MainLoop.cpp
    frontend/qt/LimeWindow.cpp
    frontend/qt/LimeView.cpp
    ${lime_MOC}
)
target_link_libraries(lime ${QT_LIBRARIES} ${PYTHON_LIBRARIES} ${Boost_LIBRARIES})
