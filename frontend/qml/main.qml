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
            MenuItem {
                text: "Quit"
                onTriggered: Qt.quit()
            }
        }
    }
    Item {
        anchors.fill: parent
        SplitView {
            anchors.fill: parent
            orientation: Qt.Vertical
            TabView {
                Layout.fillHeight: true
                id: tabs
                style: TabViewStyle {
                    frameOverlap: 0
                    tab: Item {
                        implicitWidth: 200
                        implicitHeight: 29
                        Rectangle {
                            id: tab_content
                            color: styleData.selected ? frontend.defaultBg() : "#3D3D3A"
                            visible: false
                            implicitWidth: 200
                            implicitHeight: 29
                            Text {
                                id: tab_title
                                anchors.centerIn: parent
                                text: styleData.title
                                color: frontend.defaultFg()
                                anchors.verticalCenterOffset: 1
                            }
                            Rectangle {
                                x: 0; y: 28
                                width: 200; height: 1
                                color: styleData.selected ? "transparent" : frontend.defaultBg()
                            }
                        }
                        Rectangle {
                            id: tab_mask
                            color: "transparent"
                            anchors.fill: tab_content
                            visible: false
                            Rectangle {
                                x: 21; y: 4
                                width: 200-42; height: 25
                                color: "white"
                            }
                            Image {
                                x: 0
                                width: 21; height: 29
                                source: "graphics/tab_alpha_left.png"
                            }
                            Image {
                                x: 200-21
                                width: 21; height: 29
                                source: "graphics/tab_alpha_right.png"
                            }
                        }
                        OpacityMask {
                            id: my_tab
                            anchors.fill: tab_content
                            source: tab_content
                            maskSource: tab_mask
                        }
                        Image {
                            x: 0
                            width: 21; height: 29
                            source: styleData.selected ? "graphics/tab_active_left.png" : "graphics/tab_inactive_left.png"
                        }
                        Image {
                            x: 200 - 21
                            width: 21; height: 29
                            source: styleData.selected ? "graphics/tab_active_right.png" : "graphics/tab_inactive_right.png"
                        }
                        Image {
                            x: 21
                            width: 200 - 42; height: 29
                            fillMode: Image.TileHorizontally
                            source: styleData.selected ? "graphics/tab_active_center.png" : "graphics/tab_inactive_center.png"
                        }
                    }
                    tabBar: Rectangle {
                        color: "#161713"
                        Rectangle {
                            x: 0; y: 28
                            width: parent.width; height: 1
                            color: frontend.defaultBg()
                        }
                    }
                    tabsMovable: true
                    frame: Rectangle { color: frontend.defaultBg() }
                    tabOverlap: 20
                }
                Tab {
                    anchors.fill: parent
                    title: editor.activeWindow().activeView().buffer().fileName()
                    Item {
                        id: viewItem
                        property var myView: editor.activeWindow().activeView()
                        property var viewLines: myView.buffer().rowCol(myView.buffer().size())[0]
                        Rectangle  {
                            color: frontend.defaultBg()
                            anchors.fill: parent
                        }
                        ListView {
                            id: view
                            anchors.fill: parent
                            model: viewLines
                            delegate: Text {
                                text: frontend.formatLine(myView, index)
                                textFormat: TextEdit.RichText
                                color: "white"
                            }
                            states: [
                                State {
                                    name: "ShowBars"
                                    when: view.movingVertically || view.movingHorizontally
                                    PropertyChanges {
                                        target: verticalScrollBar
                                        opacity: 0.5
                                    }
                                },
                                State {
                                    name: "HideBars"
                                    when: !view.movingVertically && !view.movingHorizontally
                                    PropertyChanges {
                                        target: verticalScrollBar
                                        opacity: 0
                                    }
                                }
                            ]
                            Rectangle {
                                id: verticalScrollBar
                                y: view.visibleArea.yPosition * view.height
                                width: 10
                                radius: width
                                height: view.visibleArea.heightRatio * view.height
                                anchors.right: view.right
                                opacity: 0
                                Behavior on opacity { PropertyAnimation {} }
                            }
                            clip: true
                        }
                    }
                }
                Tab {
                    title: "untitled"
                }
                Tab {
                    title: "untitled"
                }
            }
            Item {
                id: consoleView
                height: 100
                property var myView: editor.console()
                property var viewLines: myView.buffer().rowCol(myView.buffer().size())[0]
                Rectangle {
                    color: frontend.defaultBg()
                    anchors.fill: parent
                }
                ListView {
                    anchors.fill: parent
                    model: consoleView.viewLines
                    delegate: Text {
                        text: frontend.formatLine(consoleView.myView, index)
                        textFormat: TextEdit.RichText
                        color: "white"
                    }
                    clip: true
                }
            }
        }
    }
}
