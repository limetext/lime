import QtQuick 2.0
import QtQuick.Dialogs 1.0

FileDialog {
    id: openDialog
    title: qsTr("Open File")
    folder: {
        if (tabs.count == 0) return "/";

        var fn = myWindow.view(tabs.currentIndex).title.text;
        var sp = (Qt.platform.os == "windows") ? "\\" : "/";
        return fn.substring(0, fn.lastIndexOf(sp)+1);
    }
    selectMultiple: true
    onAccepted: {
        var _chosenFiles= openDialog.fileUrls;
        var _urls = [];
        for (var i = 0; i < _chosenFiles.length; i++) {
            _urls[i] = _chosenFiles[i].toString();
            if(_urls[i].length >= 7 && _urls[i].slice(0, 7) == "file://") {
                _urls[i] = _urls[i].slice(7)
            }
            console.log("Choosed: " + _urls[i]);
            frontend.runCommandWithArgs("open_file", {"path" : _urls[i]});
        }
    }
    onRejected: {
        console.log("Canceled.")
    }
}
