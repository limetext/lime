import QtQuick 2.0
import QtQuick.Controls 1.0
import QtQuick.Controls.Styles 1.0
import QtQuick.Layouts 1.0
import QtGraphicalEffects 1.0

ApplicationWindow {
    id: window
    width: 800
    height: 600
    menuBar: MenuBar {
        id: menu
        Menu {
            title: "Hello"
            MenuItem { text: "World" }
        }
    }
    property var myWindow

    Item {
        anchors.fill: parent
        Keys.onPressed: {
            event.accepted = frontend.handleInput(event.key, event.modifiers)
        }
        focus: true // Focus required for Keys.onPressed
        SplitView {
            anchors.fill: parent
            orientation: Qt.Vertical
            TabView {
                Layout.fillHeight: true
                id: tabs
                style: TabViewStyle {
                    frameOverlap: 0
                    tab: Item {
                        implicitWidth: 180
                        implicitHeight: 28
                        BorderImage {
                            source: styleData.selected ? "../../3rdparty/bundles/themes/soda/Soda Dark/tab-active.png" : "../../3rdparty/bundles/themes/soda/Soda Dark/tab-inactive.png"
                            border { left: 5; top: 5; right: 5; bottom: 5 }
                            width: 180
                            height: 25
                            Text {
                                id: tab_title
                                anchors.centerIn: parent
                                text: styleData.title
                                color: frontend.defaultFg()
                                anchors.verticalCenterOffset: 1
                            }
                        }
                    }
                    tabBar: Image {
                        fillMode: Image.TileHorizontally
                        source: "../../3rdparty/bundles/themes/soda/Soda Dark/tabset-background.png"
                    }
                    tabsMovable: true
                    frame: Rectangle { color: frontend.defaultBg() }
                    tabOverlap: 5
                }
                Repeater {
                    id: rpmod
                    function tmp() {
                        var ret = myWindow ? myWindow.len : 0;
                        console.log(ret);
                        return ret;
                    }
                    model: tmp()
                    delegate: Tab {
                        title: myWindow.view(index).title()
                        LimeView { myView: myWindow.view(index) }
                    }
                }

                onCurrentIndexChanged: {
                    myWindow.back().setActiveView(myWindow.view(currentIndex).back());
                }
            }
            LimeView {
                id: consoleView
                myView: frontend.console
                height: 100
            }
        }
    }
}
