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

    function getCurrentSelection() {
        if (!myView || !myView.back()) return null;
        return myView.back().sel();
    }

    function addLine() {
        view.model.append({});
    }

    function insertLine(idx) {
        view.model.insert(idx, {});
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

    ListView {
        id: view
        model: ListModel {}
        anchors.fill: parent
        boundsBehavior: Flickable.StopAtBounds
        cacheBuffer: contentHeight
        interactive: false
        clip: true
        z: 4

        property var myView
        property bool showBars: false
        property var cursor: parent.cursor

        delegate: Rectangle {
            color: "transparent"
            width: parent.width
            height: childrenRect.height

            Text {
                property var line: !myView ? null : myView.line(index)
                font.family: viewItem.fontFace
                font.pointSize: viewItem.fontSize
                text: !line ? "" : line.text+" "
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
            property var point: new Object()

            enabled: !isMinimap
            x: 0
            y: 0
            cursorShape: parent.cursor
            propagateComposedEvents: true
            height: parent.height
            width: parent.width-verticalScrollBar.width

            function colFromMouseX(lineIndex, mouseX) {

                var line = myView.back().buffer().line(myView.back().buffer().textPoint(lineIndex, 0)),
                    lineText = myView.back().buffer().substr(line);

                dummy.text = lineText;

                // if the click was farther right than the last character of
                // the line then return the last character's column
                if(mouseX > dummy.width) {
                    return myView.back().buffer().rowCol(line.b)[1]
                }

                // fixme: why do we need this magic number?
                var OFFSET_MAGIC_NUMBER = 0.5;

                // calculate a column from a given mouse x coordinate and the line text.
                var col = Math.floor(OFFSET_MAGIC_NUMBER + lineText.length * (mouseX / dummy.width));
                if (col < 0) col = 0;

                // Trying to find closest column to clicked position
                dummy.text = "<span style=\"white-space:pre\">" + lineText.substr(0, col) + "</span>";

                var d = Math.abs(mouseX - dummy.width),
                    add = (mouseX > dummy.width) ? 1 : -1;

                while (col >= 0 && col < lineText.length && Math.abs(mouseX - dummy.width) <= d) {
                    d = Math.abs(mouseX - dummy.width);
                    col += add;
                    dummy.text = "<span style=\"white-space:pre\">" + lineText.substr(0, col) + "</span>";
                }
                col -= add;

                return col;
            }

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

            onPositionChanged: {
                var item  = view.itemAt(0, mouse.y+view.contentY),
                    index = view.indexAt(0, mouse.y+view.contentY),
                    selection = getCurrentSelection();

                if (item != null && selection != null) {
                    var col = colFromMouseX(index, mouse.x);
                    point.r = myView.back().buffer().textPoint(index, col);
                    if (point.p != null && point.p != point.r) {
                        // Remove the last region and replace it with new one
                        var r = selection.get(selection.len()-1);
                        selection.substract(r);
                        selection.add(myView.region(point.p, point.r));
                        onSelectionModified();
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
                        var col = colFromMouseX(index, mouse.x);
                        point.p = myView.back().buffer().textPoint(index, col)

                        if (!ctrl) {
                            getCurrentSelection().clear();
                        }

                        getCurrentSelection().add(myView.region(point.p, point.p))
                        onSelectionModified();
                    }
                }
            }

            onDoubleClicked: {
                if (!isMinimap) {

                    var item  = view.itemAt(0, mouse.y+view.contentY),
                        index = view.indexAt(0, mouse.y+view.contentY);

                    if (item != null) {
                        var col = colFromMouseX(index, mouse.x);
                        point.p = myView.back().buffer().textPoint(index, col)

                        if (!ctrl) {
                            getCurrentSelection().clear();
                        }

                        getCurrentSelection().add(myView.back().expandByClass(myView.region(point.p, point.p), 1|2|4|8))
                        onSelectionModified();
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
            enabled: true
            width: verticalScrollBar.width
            height: view.height
            anchors.right: parent.right
            hoverEnabled: true
            drag.target: verticalScrollBar
            drag.minimumY: 0
            drag.maximumY: view.height-verticalScrollBar.height
        }
    }

    Repeater {
        id: highlightedLines
        model: (!isMinimap && getCurrentSelection()) ? getCurrentSelection().len() : 0

        delegate: Rectangle {
            property var rowcol
            property var cursor: children[0]

            color: "#444444"
            radius: 2
            border.color: "#1c1c1c"
            height: cursor.height
            y: getYPosition(rowcol)
            z: 3

            function getYPosition(rowCol) {
                if(rowCol) {
                    return rowcol[0] * (view.contentHeight/view.count) - view.contentY;
                }
                return 0;
            }

            Text {
                color: "#F8F8F0"
                font.family: viewItem.fontFace
                font.pointSize: viewItem.fontSize
            }
        }
    }

    function onSelectionModified() {
        if (myView == undefined) return;

        var selection = getCurrentSelection(),
            backend = myView.back(),
            buf = backend.buffer(),
            of = 0; // todo: rename 'of' to something more descriptive

        highlightedLines.model = myView.regionLines();

        for(var i = 0; i < selection.len(); i++) {
            var rect = highlightedLines.itemAt(i),
                s = selection.get(i);

            if (!s || !rect) continue;

            var rowcol,
                lns = getLinesFromSelection(s, buf);

            // checks if we moved cursor forward or backward
            if (s.b <= s.a) lns.reverse();
            for(var j = 0; j < lns.length; j++) {
                rect = highlightedLines.itemAt(i+of);
                of++;
                rowcol = buf.rowCol(lns[j].a);
                rect.rowcol = rowcol;
                rect.x = getCursorOffset(lns[j].a, rowcol, rect.cursor, buf);
                rowcol = buf.rowCol(lns[j].b);
                rect.width = getCursorOffset(lns[j].b, rowcol, rect.cursor, buf) - rect.x;
            }
            of--;

            rect.cursor.x = (s.b <= s.a) ? 0 : rect.width;
            rect.cursor.opacity = 1;

            var caretStyle = myView.setting("caret_style"),
                inverseCaretState = myView.setting("inverse_caret_state");

            if (caretStyle == "underscore") {
                if (inverseCaretState) {
                    rect.cursor.text = "_";
                    if (rect.width != 0)
                        rect.cursor.x -= rect.cursor.width;
                } else {
                    rect.cursor.text = "|";
                    // Shift the cursor to the edge of the character
                    rect.cursor.x -= 4;
                }
            }
        }
        // Clearing
        for(var i = of + selection.len()+1; i < highlightedLines.count; i++) {
            var rect = highlightedLines.itemAt(i);
            if (!rect) continue;
            rect.width = 0;
        }
    }

    // getCursorOffset returns the x coordinate for the cursor.
    function getCursorOffset(cursorIndex, rowcol, cursor, buf) {

        var line = buf.line(cursorIndex),
            currentLineText  = buf.substr(line);

        // text from the beginning of the line to the given column
        var textToCursor = currentLineText.substr(0, rowcol[1]);

        cursor.textFormat = TextEdit.RichText;
        cursor.text = "<span style=\"white-space:pre\">" + textToCursor + "</span>";

        var cursorOffset = cursor.width;

        cursor.textFormat = TextEdit.PlainText;
        cursor.text = "";

        return (!cursorOffset) ? 0 : cursorOffset;
    }

    // getLinesFromSelection returns an array of lines from the given
    // selection and buffer. Works like buffer.Lines()
    //
    // note: the selection could be inverted, for example if a user starts
    // selecting from the bottom up. This makes sure that the start of
    // the selection is where the user stopped selecting.
    function getLinesFromSelection(selection, buf) {
        var lines = [];

        var safeSelection = (selection.b > selection.a) ?
                    { a: selection.a, b: selection.b }:
                    { a: selection.b, b: selection.a };

        var rowCol = {
            a: buf.rowCol(safeSelection.a),
            b: buf.rowCol(safeSelection.b)
        };

        for(var i = rowCol.a[0]; i <= rowCol.b[0]; i++) {
            var lr = buf.line(buf.textPoint(i, 0)),
                a = (i == rowCol.a[0]) ? safeSelection.a : lr.a,
                b = (i == rowCol.b[0]) ? safeSelection.b : lr.b,
                res = (b > a) ? {a: a, b: b} : {a: b, b: a};
            lines.push(res);
        }

        return lines;
    }

    Timer {
        interval: 100
        repeat: true
        running: true
        onTriggered: {
            var o = 0.5 + 0.5 * Math.sin(Date.now()*0.008);

            for (var i = 0; i < highlightedLines.count; i++) {
                var rect =  highlightedLines.itemAt(i);
                rect.cursor.opacity = o;
            }
        }
    }
}
