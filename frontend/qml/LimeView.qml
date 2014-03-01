import QtQuick 2.0

Item {
    id: viewItem
    property var myView //: frontend.window(editor.activeWindow).view(editor.activeWindow.activeView)
    Rectangle  {
        color: frontend.defaultBg()
        anchors.fill: parent
    }
    onMyViewChanged: {
        console.log("myview changed");
        view.myView = myView;
    }
    ListView {
        id: view
        property var myView
        anchors.fill: parent
        model: myView ? myView.len : 0
        delegate: Text {
            property var line: myView.line(index)
            text: line.text
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
