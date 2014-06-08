import QtQuick 2.0
import QtQuick.Controls 1.1
import QtGraphicalEffects 1.0


Item {
    id: toolTipRoot
    height: toolTipContainer.height
    width: toolTipContainer.width
    visible: mouseItem.containsMouse
    clip: false
    z: parent.parent.parent.z+100

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
            toolTipRoot.x = mouse.x;
            toolTipRoot.y = mouse.y + 5;
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
