import QtQuick 2.0
import QtQuick.Layouts 1.0

Item {
    id: viewItem
    property var myView
    property bool isMinimap: false
    property int fontSize: isMinimap ? 4 : 12
    property string fontFace: "Helvetica"
    property var cursor: Qt.IBeamCursor
    property bool ctrl: false
    function sel() {
        if (!myView || !myView.back()) return null;
        return myView.back().sel();
    }
    Rectangle  {
        color: frontend.defaultBg()
        anchors.fill: parent
    }
    onMyViewChanged: {
        if (!isMinimap) {
            view.model.clear();
            view.myView = myView;
            myView.fix(viewItem);
        }
    }
    onFontSizeChanged: {
        dummy.font.pointSize = fontSize;
    }
    function addLine() {
        view.model.append({});
    }
    function insertLine(idx) {
        view.model.insert(idx, {});
    }
    ListView {
        id: view
        property var myView
        boundsBehavior: Flickable.StopAtBounds
        anchors.fill: parent
        interactive: false
        cacheBuffer: contentHeight
        model: ListModel {}

        property bool showBars: false
        property var cursor: parent.cursor
        delegate: Rectangle {
            width: parent.width
            height: childrenRect.height
            color: "transparent"
            Text {
                property var line: myView.line(index)
                font.family: viewItem.fontFace
                font.pointSize: viewItem.fontSize
                text: line.text+" "
                textFormat: TextEdit.RichText
                color: "white"
            }
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
            enabled: !isMinimap
            property var point: new Object()
            x: 0
            y: 0
            height: parent.height
            width: parent.width-verticalScrollBar.width
            propagateComposedEvents: true
            cursorShape: parent.cursor
            Text {
                // just used to measure the text.
                // If we change an actual displayed item's text,
                // there's a risk (or is it always happening?)
                // that the backend stored text data is no longer
                // connected with that text item and hence changes
                // made backend side aren't propagated.
                id: dummy
                font.family: viewItem.fontFace
                textFormat: TextEdit.RichText
                visible: false
                Component.onCompleted: {
                    dummy.font.pointSize = viewItem.fontSize
                }
            }
            function measure(line, x) {
                var line = myView.back().buffer().line(myView.back().buffer().textPoint(line, 0)),
                    str = myView.back().buffer().substr(line);
                dummy.text = str;
                // If we are clicking out of line width return end of line column
                if(x > dummy.width) return myView.back().buffer().rowCol(line.b)[1]
                // We try to start searching from somewhere close to click position
                var col = Math.floor(0.5 + str.length * x/dummy.width);

                // Trying to find closest column to clicked position
                dummy.text = "<span style=\"white-space:pre\">" + str.substr(0, col) + "</span>";
                var d = Math.abs(x - dummy.width),
                    add = (x > dummy.width) ? 1 : -1;
                while(Math.abs(x - dummy.width) <= d) {
                    d = Math.abs(x - dummy.width);
                    col += add;
                    dummy.text = "<span style=\"white-space:pre\">" + str.substr(0, col) + "</span>";
                }
                col -= add;

                return col;
            }
            onPositionChanged: {
                var item  = view.itemAt(0, mouse.y+view.contentY),
                    index = view.indexAt(0, mouse.y+view.contentY),
                    s = sel();
                if (item != null && sel != null) {
                    var col = measure(index, mouse.x);
                    point.r = myView.back().buffer().textPoint(index, col);
                    if (point.p != null && point.p != point.r) {
                        // Remove the last region and replace it with new one
                        var r = s.get(s.len()-1);
                        s.substract(r);
                        s.add(myView.region(point.p, point.r));
                    }
                }
                point.r = null;
            }
            onPressed: {
                // TODO:
                // Changing caret position doesn't work on empty lines
                if (!isMinimap) {
                    var item  = view.itemAt(0, mouse.y+view.contentY),
                        index = view.indexAt(0, mouse.y+view.contentY);
                    if (item != null) {
                        var col = measure(index, mouse.x);
                        point.p = myView.back().buffer().textPoint(index, col)

                        if (!ctrl) sel().clear();
                        sel().add(myView.region(point.p, point.p))
                    }
                }
            }
            onDoubleClicked: {
                if (!isMinimap) {
                    var item  = view.itemAt(0, mouse.y+view.contentY),
                        index = view.indexAt(0, mouse.y+view.contentY);
                    if (item != null) {
                        var col = measure(index, mouse.x);
                        point.p = myView.back().buffer().textPoint(index, col)

                        if (!ctrl) sel().clear();
                        sel().add(myView.back().expandByClass(myView.region(point.p, point.p), 1|2|4|8))
                    }
                }
            }
            onWheel: {
                var delta = wheel.pixelDelta,
                    scaleFactor = 30;
                if (delta.x == 0 && delta.y == 0) {
                    delta = wheel.angleDelta;
                    scaleFactor = 15;
                }
                view.flick(delta.x*scaleFactor, delta.y*scaleFactor);
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
        z: 4
    }
    Repeater {
        id: regs
        model: (!isMinimap && sel()) ? sel().len() : 0
        delegate: Rectangle {
            property var rowcol
            property var cursor: children[0]
            color: "#444444"
            radius: 2
            border.color: "#1c1c1c"
            height: cursor.height
            Text {
                color: "#F8F8F0"
                font.family: viewItem.fontFace
                font.pointSize: viewItem.fontSize
            }
            y: rowcol ? rowcol[0]*(view.contentHeight/view.count)-view.contentY : 0;
            z: 3
        }
    }
    Timer {
        interval: 100
        running: true
        repeat: true
        function measure(p, rowcol, cursor, buf) {
            var line = buf.line(p),
                str  = buf.substr(line);

            str = str.substr(0, rowcol[1]);
            cursor.textFormat = TextEdit.RichText;
            cursor.text = "<span style=\"white-space:pre\">" + str + "</span>";
            var ret = cursor.width;
            cursor.textFormat = TextEdit.PlainText;
            cursor.text = "";

            return (!ret) ? 0 : ret;
        }
        // Works like buffer.Lines()
        function lines(sel, buf) {
            var lines = new Array(),
                ss = (sel.b > sel.a) ? {a: sel.a, b: sel.b} : {a: sel.b, b: sel.a},
                rc = {a: buf.rowCol(ss.a), b: buf.rowCol(ss.b)};

            for(var i = rc.a[0]; i <= rc.b[0]; i++) {
                var lr = buf.line(buf.textPoint(i, 0)),
                    a = (i == rc.a[0]) ? ss.a : lr.a,
                    b = (i == rc.b[0]) ? ss.b : lr.b,
                    res = (b > a) ? {a: a, b: b} : {a: b, b: a};
                lines.push(res);
            }
            return lines;
        }
        onTriggered: {
            if (myView == undefined) return;
            var s = sel(),
                back = myView.back(),
                buf = back.buffer(),
                of = 0;
            regs.model = myView.regionLines();
            for(var i = 0; i < s.len(); i++) {
                var rect = regs.itemAt(i),
                    mysel = s.get(i);
                if (!mysel || !rect) continue;

                var rowcol,
                    lns = lines(mysel, buf);

                if (mysel.b <= mysel.a) lns.reverse();
                for(var j = 0; j < lns.length; j++) {
                    rect = regs.itemAt(i+of);
                    of++;
                    rowcol = buf.rowCol(lns[j].a);
                    rect.rowcol = rowcol;
                    rect.x = measure(lns[j].a, rowcol, rect.cursor, buf);
                    rowcol = buf.rowCol(lns[j].b);
                    rect.width = measure(lns[j].b, rowcol, rect.cursor, buf) - rect.x;
                }
                of--;

                rect.cursor.x = (mysel.b <= mysel.a) ? 0 : rect.width;
                rect.cursor.opacity = 0.5 + 0.5 * Math.sin(Date.now()*0.008);;

                var style = myView.setting("caret_style"),
                    inv = myView.setting("inverse_caret_state");
                if (style == "underscore") {
                    if (inv) {
                        rect.cursor.text = "_";
                    } else {
                        rect.cursor.text = "|";
                        // Shift the cursor to the edge of the character
                        rect.cursor.x -= 4;
                    }
                }
            }
            // Clearing
            for(var i = of+s.len()+1; i < regs.count; i++) {
                var rect = regs.itemAt(i);
                if (!rect) continue;
                rect.width = 0;
            }
        }
    }
}
