import QtQuick 2.0
import QtQuick.Controls 1.0


ApplicationWindow {
     id: window
    objectName: "rc"
    property var myProperty
    width: 100
    height: 40


    function fucn1(birthday) {
        var ageDifMs = Date.now() - birthday.getTime();
          var ageDate = new Date(ageDifMs); // miliseconds from epoch
          return Math.abs(ageDate.getUTCFullYear() - 1970);


    }



    function functionShouldReturn_10() {
        return 10;

    }


    function functionShouldExpandWindow() {

            window.width = 101 + Math.random();


    }


    function functionShouldCallmodelMyFunction() {
        myvar.myFunction();
    }

    function functionShouldCallmodelMyFunctionWithArgument() {
        myvar.myFunctionWithInt(54);


    }



}


