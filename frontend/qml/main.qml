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
    property bool ctrl

    menuBar: MenuBar {
        id: menu
        Menu {
            title: qsTr("&File")
            MenuItem {
                text: qsTr("&New File")
                shortcut: "Ctrl+N"
                onTriggered: frontend.runCommand("new_file");
            }
            MenuItem {
                text: qsTr("&Open File...")
                shortcut: "Ctrl+O"
                onTriggered: frontend.runCommand("open_file");
            }
            MenuItem {
                text: qsTr("&Save")
                shortcut: "Ctrl+S"
                onTriggered: frontend.runCommand("save");
            }
            MenuItem {
                text: qsTr("&Save As...")
                shortcut: "Shift+Ctrl+S"
                // TODO(.) : qml doesn't have a ready dialog like FileDialog
                // onTriggered: saveAsDialog.open()
            }
            MenuSeparator{}
            MenuItem {
                text: qsTr("&New Window")
                shortcut: "Shift+Ctrl+N"
                onTriggered: frontend.runCommand("new_window");
            }
            MenuItem {
                text: qsTr("&Close Window")
                shortcut: "Shift+Ctrl+W"
                onTriggered: frontend.runCommand("close_window");
            }
            MenuSeparator{}
            MenuItem {
                text: qsTr("&Close File")
                shortcut: "Ctrl+W"
                onTriggered: frontend.runCommand("close");
            }
            MenuItem {
                text: qsTr("&Close All Files")
                onTriggered: frontend.runCommand("close_all");
            }
            MenuSeparator{}
            MenuItem {
                text: qsTr("&Quit")
                shortcut: "Ctrl+Q"
                onTriggered: Qt.quit(); // frontend.runCommand("quit");
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
            ctrl = (event.modifiers && Qt.ControlModifier) ? true : false;
            event.accepted = frontend.handleInput(event.key, event.modifiers)
        }
        Keys.onReleased: {
            ctrl = false;
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
                        onItemAdded: {
                            if (myWindow.activeViewIndex() < tabs.count)
                                tabs.currentIndex = myWindow.activeViewIndex()
                        }
                    }
                    function resetminimap() {
                        // TODO(.): This is conflicts on new file on new file the
                        //          active_view should be the new file but it changes to first tab
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

                    function scroll() {
                        var p = percentage(realView);
                        children[1].contentY = p*(children[1].contentHeight-height);
                        if (!ma.drag.active) {
                            minimapArea.y =  p*(height-minimapArea.height)
                        }
                    }

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
                        width: parent.width
                        height: parent.realView ? parent.realView.visibleArea.heightRatio*parent.children[1].contentHeight : parent.height
                        color: "white"
                        opacity: 0.1
                        onYChanged: {
                            if (ma.drag.active) {
                                parent.realView.contentY = y*(parent.realView.contentHeight-parent.realView.height)/(parent.height-height);
                            }
                        }
                        onHeightChanged: {
                            parent.scroll();
                        }
                        MouseArea {
                            id: ma
                            drag.target: parent
                            anchors.fill: parent
                            drag.minimumX: 0
                            drag.minimumY: 0
                            drag.maximumY: parent.parent.height-height
                            drag.maximumX: parent.parent.width-width
                        }
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
        title: qsTr("Open File")
        // TODO(.) : folder should be set to current view directory
        // folder: myWindow.view(tabs.currentIndex).title.text
        // TODO(.) : Selecting multiple files should be enabled
        // selectMultiple: true
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
