import QtQuick 2.0
import QtQuick.Layouts 1.0

Item {
    id: viewItem
    property var myView
    property bool isMinimap: false
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
        property bool showBars: false
        delegate: Text {
            property var line: myView.line(index)
            font.pointSize: isMinimap ? 4 : 12
            text: line.text
            textFormat: TextEdit.RichText
            color: "white"
        }
        states: [
            State {
                name: "ShowBars"
                when: view.movingVertically || view.movingHorizontally || view.flickingVertically || view.flickingHorizontally
                PropertyChanges {
                    target: view
                    showBars: true
                }
            },
            State {
                name: "HideBars"
                when: !view.movingVertically && !view.movingHorizontally && !view.flickingVertically && !view.flickingHorizontally
                PropertyChanges {
                    target: view
                    showBars: false
                }
            }
        ]
        MouseArea {
            x: 0
            y: 0
            height: parent.height
            width: parent.width-verticalScrollBar.width
            propagateComposedEvents: true
            onWheel: {
                view.flick(0, wheel.angleDelta.y*100);
                wheel.accepted = true;
            }
        }

        Rectangle {
            id: verticalScrollBar
            visible: !isMinimap
            width: 10
            radius: width
            height: view.visibleArea.heightRatio * view.height
            anchors.right: view.right
            opacity: (view.showBars || ma.containsMouse || ma.drag.active) ? 0.5 : 0.0
            onYChanged: {
                if (ma.drag.active) {
                    view.contentY = y*(view.contentHeight-view.height)/(view.height-height);
                }
            }
            states: [
                State {
                    when: !ma.drag.active
                    PropertyChanges {
                        target: verticalScrollBar
                        y: view.visibleArea.yPosition*view.height
                    }
                }
            ]
            Behavior on opacity { PropertyAnimation {} }
        }
        MouseArea {
            id: ma
            width: verticalScrollBar.width
            height: view.height
            anchors.right: parent.right
            hoverEnabled: true
            drag.target: verticalScrollBar
            drag.minimumY: 0
            drag.maximumY: view.height-verticalScrollBar.height
            enabled: true
        }

        clip: true
    }
}
