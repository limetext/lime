import QtQuick 2.0
import QtQuick.Controls 1.0
import QtQuick.Controls.Styles 1.0
import QtQuick.Dialogs 1.0
import QtQuick.Layouts 1.0
import QtGraphicalEffects 1.0

ApplicationWindow {
    id: window
    width: 800
    height: 600

    property var myWindow

    menuBar: MenuBar {
        id: menu
        Menu {
            title: qsTr("&File")
            MenuItem {
                text: qsTr("&New")
                shortcut: "Ctrl+N"
                onTriggered: editor.newWindow()
            }
            MenuItem {
                text: qsTr("&Open")
                shortcut: "Ctrl+O"
                onTriggered: openDialog.open()
            }
        }
    }

    statusBar: StatusBar {
        id: statusBar
        style: StatusBarStyle {
            background: Image {
               source: "../../3rdparty/bundles/themes/soda/Soda Dark/status-bar-background.png"
            }
        }

        property color textColor: "#969696"

        RowLayout {
            anchors.fill: parent
            id: statusBarRowLayout
            spacing: 15

            RowLayout {
                anchors.fill: parent
                spacing: 3

                Label {
                    text: "git branch: master"
                    color: statusBar.textColor
                }

                Label {
                    text: "INSERT MODE"
                    color: statusBar.textColor
                }

                Label {
                    id: statusBarCaretPos
                    text: "Line xx, Column yy"
                    color: statusBar.textColor
                }
            }

            Label {
                id: statusBarIndent
                text: "Tab Size/Spaces: 4"
                color: statusBar.textColor
                Layout.alignment: Qt.AlignRight
            }

            Label {
                id: statusBarLanguage
                text: "Go"
                color: statusBar.textColor
                Layout.alignment: Qt.AlignRight
            }
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
            SplitView {
                Layout.fillHeight: true
                TabView {
                    Layout.fillHeight: true
                    Layout.fillWidth: true
                    id: tabs
                    style: TabViewStyle {
                        frameOverlap: 0
                        tab: Item {
                            implicitWidth: 180
                            implicitHeight: 28
                             ToolTip {
                                id: tooltip
                                backgroundColor: "#BECCCC66"
                                textColor: "black"
                                font.pointSize: 8
                                text: styleData.title
                            }
                            BorderImage {
                                source: styleData.selected ? "../../3rdparty/bundles/themes/soda/Soda Dark/tab-active.png" : "../../3rdparty/bundles/themes/soda/Soda Dark/tab-inactive.png"
                                border { left: 5; top: 5; right: 5; bottom: 5 }
                                width: 180
                                height: 25
                                Text {
                                    id: tab_title
                                    anchors.centerIn: parent
                                    text: styleData.title.replace(/^.*[\\\/]/, '')
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
                            title: myWindow.view(index).title.text
                            LimeView {
                                id: view
                                myView: myWindow.view(index)
                                Timer {
                                    // TODO(.): Why is view null sometimes?
                                    //          Is it just a variant of https://github.com/go-qml/qml/issues/68?
                                    interval: 100
                                    running: view.myView == null
                                    repeat: true
                                    onTriggered: {
                                        view.myView = myWindow.view(index);
                                        if (index == tabs.currentIndex) {
                                            tabs.resetminimap();
                                        }
                                    }
                                }
                            }
                        }
                    }
                    function resetminimap() {
                        myWindow.back().setActiveView(myWindow.view(currentIndex).back());

                        minimap.myView = myWindow.view(currentIndex);
                        minimap.realView = tabs.getTab(currentIndex).item.children[1];
                    }
                    Component.onCompleted: {
                        resetminimap();
                    }
                    onCurrentIndexChanged: {
                        resetminimap();
                    }
                }
                LimeView {
                    Layout.maximumWidth: 100
                    Layout.minimumWidth: 100
                    Layout.preferredWidth: 100
                    width: 100
                    isMinimap: true
                    cursor: Qt.ArrowCursor
                    property var realView
                    property var oldView

                    function scroll() { children[1].contentY = percentage(realView)*(children[1].contentHeight-height); }

                    onRealViewChanged: {
                        if (oldView) {
                            oldView.contentYChanged.disconnect(scroll);
                        }
                        realView.contentYChanged.connect(scroll);
                        oldView = realView;
                    }
                    function percentage(view) {
                        return view.visibleArea.yPosition/(1-view.visibleArea.heightRatio);
                    }
                    id: minimap
                    Rectangle {
                        id: minimapArea
                        y: parent.percentage(parent.children[1])*(parent.height-height)
                        width: parent.width
                        height: parent.realView ? parent.realView.visibleArea.heightRatio*parent.children[1].contentHeight : parent.height
                        color: "white"
                        opacity: 0.1
                    }
                }
            }
            LimeView {
                id: consoleView
                myView: frontend.console
                Timer {
                    // TODO(.): Why is view null sometimes?
                    //          Is it just a variant of https://github.com/go-qml/qml/issues/68?
                    interval: 100
                    running: consoleView.myView == null
                    repeat: true
                    onTriggered: {
                        consoleView.myView = myWindow.view(index);
                    }
                }

                height: 100
            }
        }
    }

    FileDialog {
        id: openDialog
        title: qsTr("Please choose a file:")
        onAccepted: {
            var _url = openDialog.fileUrl.toString()
            if(_url.length >= 7 && _url.slice(0, 7) == "file://") {
                _url = _url.slice(7)
            }
            console.log("Choosed: " + _url);
            myWindow.back().openFile(_url, 0);
        }
        onRejected: {
            console.log("Canceled.")
        }
    }
}
