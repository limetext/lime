import QtQuick 2.0
import QtQuick.Dialogs 1.0

FileDialog {
        id: openDialog
        title: qsTr("Open File")
        // TODO(.) : folder should be set to current view directory
        // folder: myWindow.view(tabs.currentIndex).title.text
        // TODO(.) : Selecting multiple files should be enabled
        // selectMultiple: true
        onAccepted: {
            var _url = openDialog.fileUrl.toString()
            if(_url.length >= 7 && _url.slice(0, 7) == "file://") {
                _url = _url.slice(7)
            }
            console.log("Choosed: " + _url);
            frontend.runCommandWithArgs("open_file", {"path" : _url});
        }
        onRejected: {
            console.log("Canceled.")
        }
    }