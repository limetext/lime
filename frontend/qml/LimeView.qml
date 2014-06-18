import QtQuick 2.0
import QtQuick.Layouts 1.0

Item {
    id: viewItem
    property var myView
    property bool isMinimap: false
    property double fontSize: isMinimap ? 4 : 12
    property string fontFace: "Helvatica"
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
        if (myView != null) {
            viewItem.fontSize = isMinimap ? 4 : parseFloat(myView.setting("font_size"));
            viewItem.fontFace = String(myView.setting("font_face"));
        }
    }
    ListView {
        id: view
        property var myView
        boundsBehavior: Flickable.StopAtBounds
        anchors.fill: parent
        interactive: false
        cacheBuffer: contentHeight
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
                font.pointSize: viewItem.fontSize
                textFormat: TextEdit.RichText
                visible: false
            }
            function measure(el, line, mouse) {
                var line = myView.back().buffer().line(myView.back().buffer().textPoint(line, 0));
                // If we are clicking out of line width return end of line column
                if(mouse.x > el.width) return myView.back().buffer().rowCol(line.b)[1]
                var str  = myView.back().buffer().substr(line);
                // We try to start searching from somewhere close to click position
                var col  = Math.floor(0.5 + str.length * mouse.x/el.width);

                // Trying to find closest column to clicked position
                dummy.text = "<span style=\"white-space:pre\">" + str.substr(0, col) + "</span>";
                var d = Math.abs(mouse.x - dummy.width)
                var add = (mouse.x > dummy.width) ? 1 : -1
                while(Math.abs(mouse.x - dummy.width) <= d) {
                    d = dummy.width - mouse.x;
                    col += add;
                    dummy.text = "<span style=\"white-space:pre\">" + str.substr(0, col) + "</span>";
                }
                col -= add;

                return col;
            }
            onPositionChanged: {
                var item  = view.itemAt(0, mouse.y);
                var index = view.indexAt(0, mouse.y);
                var s = sel();
                if (item != null && sel != null) {
                    var col = measure(item, index, mouse);
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
                // Multi cursor on holding ctrl key
                if (!isMinimap) {
                    var item  = view.itemAt(0, mouse.y)
                    var index = view.indexAt(0, mouse.y)
                    if (item != null) {
                        var col = measure(item, index, mouse)
                        point.p = myView.back().buffer().textPoint(index, col)
                        // If ctrl is not pressed clear the regions
                        if (!false) sel().clear()
                        sel().add(myView.region(point.p, point.p))
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
        id: regs
        model: (!isMinimap && sel()) ? sel().len() : 0
        Rectangle {
            property var rowcol
            property var cursor: children[0]
            color: "white"
            radius: 2
            opacity: 0.6
            height: cursor.height
            Text {
                color: "white"
                font.family: viewItem.fontFace
                font.pointSize: fontSize
            }
            y: rowcol ? rowcol[0]*(view.contentHeight/view.count)-view.contentY : 0;
        }
    }
    Timer {
        interval: 100
        running: true
        repeat: true
        function measure(p, rowcol, cursor, buf) {
            var line = buf.line(p);
            if (!line) return 0;
            var str  = buf.substr(line);
            if (!str) return 0;

            str = str.substr(0, rowcol[1]);
            cursor.textFormat = TextEdit.RichText;
            cursor.text = "<span style=\"white-space:pre\">" + str + "</span>";
            var ret = cursor.width;
            cursor.textFormat = TextEdit.PlainText;
            cursor.text = "";

            return (ret == null) ? 0 : ret;
        }
        // Works like buffer.Lines()
        function lines(sel, buf) {
            if (!sel) return;
            var lines  = new Array();
            var sel = (sel.b > sel.a) ? {a: sel.a, b: sel.b} : {a: sel.b, b: sel.a};
            var rc = {a: buf.rowCol(sel.a), b: buf.rowCol(sel.b)};

            for(var i = rc.a[0]; i <= rc.b[0]; i++) {
                var mysel = buf.line(buf.textPoint(i, 0));
                // Sometimes null don't know why
                if (!mysel) continue;
                var a = (i == rc.a[0]) ? sel.a : mysel.a;
                var b = (i == rc.b[0]) ? sel.b : mysel.b;
                var res = (b > a) ? {a: a, b: b} : {a: b, b: a};
                lines.push(res);
            }
            return lines;
        }
        onTriggered: {
            // TODO(.): not too happy about actively polling like this
            var s = sel();
            if (!s) return;
            var back = myView.back()
            if (!back) return;
            var buf = back.buffer();
            if (!buf) return;
            var of = 0;
            regs.model = myView.lines();
            for(var i = 0; i < s.len(); i++) {
                var rect = regs.itemAt(i);
                var mysel = s.get(i);
                if (!mysel || !rect) continue;

                var rowcol;
                var lns = lines(mysel, buf);

                if (mysel.b <= mysel.a) lns.reverse();
                for(var j = 0; j < lns.length; j++,of++) {
                    rect = regs.itemAt(i+of);
                    if (!rect) continue;
                    rowcol = buf.rowCol(lns[j].a);
                    rect.rowcol = rowcol;
                    rect.x = measure(lns[j].a, rowcol, rect.cursor, buf);
                    rowcol = buf.rowCol(lns[j].b);
                    var tmp = measure(lns[j].b, rowcol, rect.cursor, buf)
                    rect.width = tmp - rect.x;
                }

                if (!rect) continue;
                rect.cursor.x = (mysel.b <= mysel.a) ? -3 : rect.width;
                rect.cursor.opacity = 0.5 + 0.5 * Math.sin(Date.now()*0.008);;

                var style = myView.setting("caret_style");
                var inv = myView.setting("inverse_caret_state");
                if (style == "underscore") {
                    if (inv) {
                        rect.cursor.text = "_";
                    } else {
                        rect.cursor.text = "|";
                        rect.x -= 2;
                    }
                }
            }
            // Clearing
            for(var i = of+s.len()+1; i < regs.count; i++) {
                var rect = regs.itemAt(i);
                if (rect == null) continue;
                rect.width = 0;
            }
        }
    }
}
