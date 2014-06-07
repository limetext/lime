import QtQuick 2.0
import QtQuick.Controls 1.1
import QtGraphicalEffects 1.0
 
 
Item {
    id: toolTipRoot
    height: toolTipContainer.height
    visible: false
    clip: false
    z: 999999999
 
    property alias text: toolTip.text
    property alias backgroundColor: content.color
    property alias textColor: toolTip.color
    property alias font: toolTip.font
 
    function onMouseHover(x, y)
    {
        toolTipRoot.x = x;
        toolTipRoot.y = y + 5;
    }
 
    function onVisibleStatus(flag)
    {
        toolTipRoot.visible = flag;
    }
 
    Component.onCompleted: {
        var newObject = Qt.createQmlObject('import QtQuick 2.0; MouseArea {signal mouserHover(int x, int y); signal showChanged(bool flag); anchors.fill:parent; hoverEnabled: true; onPositionChanged: {mouserHover(mouseX, mouseY)} onEntered: {showChanged(true)} onExited:{showChanged(false)}}',
            toolTipRoot.parent, "mouseItem");
        newObject.mouserHover.connect(onMouseHover);
        newObject.showChanged.connect(onVisibleStatus);
    }
 
    Item {
        id: toolTipContainer
        width: content.width + toolTipShadow.radius
        height: content.height + toolTipShadow.radius
 
        Rectangle {
            id: content
            width: toolTipRoot.width
            height: toolTip.contentHeight + 10
 
            Text {
                id: toolTip
                anchors {fill: parent; margins: 5}
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
