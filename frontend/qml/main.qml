import QtQuick 2.0
import QtQuick.Controls 1.0
import QtQuick.Controls.Styles 1.0

ApplicationWindow {
    id: window;
    width: 800; height: 600;
    menuBar: MenuBar {
        id: menu;
        Menu {
            title: "Hello"
            MenuItem { text: "World" }
        }
    }
    SplitView {
        anchors.fill: parent;
        TabView {
            id: tabs;
            anchors.fill: parent;
            Tab {
                anchors.fill: parent;
                title: editor.activeWindow().activeView().buffer().fileName();
                Item {
                    id: viewItem;
                    property var myView: editor.activeWindow().activeView()
                    property var viewLines: myView.buffer().rowCol(myView.buffer().size())[0]
                    Rectangle  {
                        color: frontend.defaultBg();
                        anchors.fill: parent
                    }
                    ListView {
                        id: view;
                        anchors.fill: parent;
                        model: viewLines
                        delegate: Text {
                            text: frontend.formatLine(myView, index)
                            textFormat: TextEdit.RichText
                            color: "white"
                        }
                        states: State {
                            name: "ShowBars"
                            when: view.movingVertically || view.movingHorizontally
                            PropertyChanges { target: verticalScrollBar; opacity: 0.5 }
                        }
                        Rectangle {
                            id: verticalScrollBar
                            y: view.visibleArea.yPosition * view.height;
                            width: 10;
                            radius: width;
                            height: view.visibleArea.heightRatio * view.height
                            anchors.right: view.right
                            opacity: 0
                            Behavior on opacity { PropertyAnimation {} }
                        }
                    }
                }
            }
        }
    }
}
