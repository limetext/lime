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
                Tab {
                    anchors.fill: parent
                    title: editor.activeWindow.activeView.buffer().filename
                    Item {
                        id: viewItem
                        property var myView: editor.activeWindow.activeView
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
