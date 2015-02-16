import QtQuick 2.0
import QtQuick.Controls 1.1
import QtGraphicalEffects 1.0


Item {
    id: toolTipRoot
    height: toolTipContainer.height
    width: toolTipContainer.width
    visible: false
    clip: false
    z: parent.parent.parent.z+100
    opacity: visible ? 1 : 0
    Behavior on opacity { PropertyAnimation { duration: 250} }

    property alias text: toolTip.text
    property alias backgroundColor: content.color
    property alias textColor: toolTip.color
    property alias font: toolTip.font


    MouseArea {
        id: mouseItem
        anchors.fill: parent;
        hoverEnabled: true;
        acceptedButtons: Qt.NoButton
        onPositionChanged: {
            function getAbsolutePosition(node) {
                var returnPos = {};
                returnPos.x = 0;
                returnPos.y = 0;
                if(node !== undefined && node !== null) {
                    var parentValue = getAbsolutePosition(node.parent);
                    returnPos.x = parentValue.x + node.x;
                    returnPos.y = parentValue.y + node.y;
                }
                return returnPos;
            }
            var pos = getAbsolutePosition(this);
            pos.x += mouse.x;
            pos.y += mouse.y;
            toolTipContainer.x = pos.x;
            toolTipContainer.y = pos.y + 5;
        }
        Timer {
            interval: 500
            running: mouseItem.containsMouse
            repeat: false
            onTriggered: {
                toolTipRoot.visible = true;
            }
        }
        onExited: {
            toolTipRoot.visible = false;
        }
    }

    Component.onCompleted: {
        mouseItem.parent = toolTipRoot.parent;
    }

    Item {
        id: toolTipContainer
        width: content.width + toolTipShadow.radius
        height: content.height + toolTipShadow.radius
        z: toolTipRoot.z

        Rectangle {
            id: content
            width: toolTip.width + 10
            height: toolTip.contentHeight + 10
            Text {
                x: 5
                y: 5
                id: toolTip
                wrapMode: Text.WrapAnywhere
            }
        }
    }

    DropShadow {
        id: toolTipShadow
        z: toolTipRoot.z
        anchors.fill: source
        cached: true
        horizontalOffset: 4
        verticalOffset: 4
        radius: 8.0
        samples: 16
        color: "#80000000"
        smooth: true
        source: toolTipContainer
    }
}
