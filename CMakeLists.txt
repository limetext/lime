cmake_minimum_required(VERSION 2.8.9)

# use for example -DQt5Widgets_DIR=/usr/local/Qt-5.0.0/lib/cmake/Qt5Widgets/
find_package(Qt5Widgets)
find_package(Boost COMPONENTS Python REQUIRED)
find_package(PythonLibs REQUIRED)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

if(NOT Qt5Widgets_FOUND)
    find_package(Qt4 COMPONENTS QtGui QtCore REQUIRED)
    include(${QT_USE_FILE})
    add_definitions(${QT_DEFINITIONS})
    qt4_wrap_cpp(lime_MOC frontend/qt/MainLoop.h frontend/qt/LimeWindow.h)
else()
    set(CMAKE_AUTOMOC ON)
    set(CMAKE_INCLUDE_CURRENT_DIR ON)
endif()

include_directories(${PYTHON_INCLUDE_DIRS} ${Boost_INCLUDE_DIRS})

add_executable(lime MACOSX_BUNDLE
    frontend/qt/LimeMinimap.cpp
    frontend/qt/LimeView.cpp
    frontend/qt/LimeViewWidget.cpp
    frontend/qt/LimeWindow.cpp
    frontend/qt/main.cpp
    frontend/qt/MainLoop.cpp
    ${lime_MOC}
)
target_link_libraries(lime ${QT_LIBRARIES} ${PYTHON_LIBRARIES} ${Boost_LIBRARIES})

if(Qt5Widgets_FOUND)
    qt5_use_modules(lime Widgets)
endif()

if(APPLE)
    set_target_properties(lime
        PROPERTIES
        MACOSX_BUNDLE_INFO_PLIST ${PROJECT_SOURCE_DIR}/Info.plist
    )
endif()
