import QtQuick 2.0

Item {
    width: 800; height: 600;
     Rectangle  {
        width: parent.width; height: parent.height; color: frontend.defaultBg();
        anchors.bottom: parent.bottom
    }
    ListView {
        width: parent.width;
        height: parent.height;
        model: lines.len
        delegate: Text {
            text: lines.formatLine(index)
            textFormat: TextEdit.RichText
        }
    }
}
