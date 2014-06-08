import QtQuick 2.0
import QtQuick.Layouts 1.0

Item {
    id: viewItem
    property var myView
    property bool isMinimap: false
    property double fontSize: isMinimap ? 4 : parseFloat(myView.setting("font_size"))
    property string fontFace: String(myView.setting("font_face"))
    property var cursor: Qt.IBeamCursor
    function sel() {
        if (!myView || !myView.back()) return null;
        return myView.back().sel();
    }
    Rectangle  {
        color: frontend.defaultBg()
        anchors.fill: parent
    }
    onMyViewChanged: {
        view.myView = myView;
    }
    ListView {
        id: view
        property var myView
        anchors.fill: parent
        onMyViewChanged: {
            if (myView != null) {
                model = myView.len;
            }
            console.log(myView);
        }

        property bool showBars: false
        property var cursor: parent.cursor
        delegate: Text {
            property var line: myView.line(index)
            font.family: viewItem.fontFace
            font.pointSize: viewItem.fontSize
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
            cursorShape: parent.cursor
            function measure(el, line, mouse) {
                var line = myView.back().buffer().line(myView.back().buffer().textPoint(line, 0));
                // If we are clicking out of line width return end of line column
                if(mouse.x > el.width) return myView.back().buffer().rowCol(line.b)[1]
                var str  = myView.back().buffer().substr(line);
                // We try to start searching from somewhere close to click position
                var col  = Math.floor(0.5 + str.length * mouse.x/el.width);

                // Trying to find closest column to clicked position
                el.text = "<span style=\"white-space:pre\">" + str.substr(0, col) + "</span>";
                var d = Math.abs(mouse.x - el.width)
                var add = (mouse.x > el.width) ? 1 : -1
                while(Math.abs(mouse.x - el.width) <= d) {
                    d = el.width - mouse.x
                    col += add
                    el.text = "<span style=\"white-space:pre\">" + str.substr(0, col) + "</span>";
                }
                col -= add

                el.text = el.line.text;
                return col
            }
            onClicked: {
                // TODO: If ctrl key is pressed or we are selecting text
                // we should do sth else
                if (!isMinimap) {
                    var item  = view.itemAt(0, mouse.y)
                    var index = view.indexAt(0, mouse.y)
                    if (item != null) {
                        var col   = measure(item, index, mouse)
                        myView.back().sel().clear()
                        myView.addR(myView.back().buffer().textPoint(index, col))
                    }
                }
            }
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
    Repeater {
        model: (!isMinimap && sel()) ? sel().len() : 0
        delegate: Text {
            property var rowcol
            Timer {
                interval: 100
                running: true
                repeat: true
                function measure(mysel, rowcol) {
                    var back = myView.back();
                    if (!back) { return 0; }
                    var buf = back.buffer();
                    if (!buf) { return 0; }
                    var line = buf.line(mysel.b);
                    // TODO(.): would be better to use proper font metrics
                    // TODO(.): This assignment makes qml panic with the confusing error
                    // "panic: cannot use int as a int"
                    // line.b = mysel.b;
                    var str = buf.substr(line);
                    str = str.substr(0, rowcol[1]);
                    parent.textFormat = TextEdit.RichText;
                    var old = parent.text;
                    parent.text = "<span style=\"white-space:pre\">" + str + "</span>";
                    var ret = parent.width;
                    parent.textFormat = TextEdit.PlainText;
                    parent.text = old;

                    if (ret == null) {
                        ret = 0;
                    }
                    return ret;
                }
                onTriggered: {
                    // TODO(.): not too happy about actively polling like this
                    var s = sel();
                    if (!s) return;
                    if (index >= s.len()) {
                        return;
                    }
                    var mysel = s.get(index);
                    parent.rowcol = myView.back().buffer().rowCol(mysel.b);
                    var width = measure(mysel, parent.rowcol);

                    parent.x = width;
                    parent.opacity = 0.5 + 0.5 * Math.sin(Date.now()*0.008);

                    var style = myView.setting("caret_style");
                    var inv = myView.setting("inverse_caret_state");
                    if (style == "underscore") {
                        if (inv) {
                            text = "_";
                        } else {
                            text = "|";
                            parent.x -= 2;
                        }
                    }
                }

            }
            y: rowcol ? rowcol[0]*(view.contentHeight/view.count)-view.contentY : 0;
            font.family: viewItem.fontFace
            font.pointSize: fontSize
            color: "white"
        }
    }

}
