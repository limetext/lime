cmake_minimum_required(VERSION 2.6)

find_package(Qt4 COMPONENTS QtGui QtCore REQUIRED)
include(${QT_USE_FILE})
find_package(PythonLibs REQUIRED)
add_definitions(${QT_DEFINITIONS})
include_directories(${PYTHON_INCLUDE_DIRS})

qt4_wrap_cpp(lime_MOC frontend/qt/mainloop.h)
add_executable(lime MACOSX_BUNDLE frontend/qt/main.cpp ${lime_MOC})
target_link_libraries(lime ${QT_LIBRARIES} ${PYTHON_LIBRARIES})
