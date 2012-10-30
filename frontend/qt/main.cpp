#include <QApplication>
#include <Python.h>
#include <math.h>
#include "LimeWindow.h"
#include "MainLoop.h"

static void new_window(object arg)
{
    boost::shared_ptr<LimeWindow> w(new LimeWindow(arg));
    setattr(arg, "qtLimeWindow", w);

    w->resize(600, 400);
    w->show();
}

int main(int argc, char** argv)
{
    QApplication app(argc, argv);
    Py_Initialize();
    int ret = -1;
    {
        MainLoop *ml = MainLoop::GetInstance();
        try
        {
            register_ptr_to_python<boost::shared_ptr<LimeWindow> >();

            object editor = ml->GetEditor();

            editor.attr("new_window_event") += new_window;
            object window = editor.attr("new_window")();
            window.attr("open_file")("../backend/backend.py");
            window.attr("open_file")("../frontend/qt/lime.py");
        }
        catch(boost::python::error_already_set& e)
        {
            PyErr_Print();
            PyErr_Clear();
        }

        ret = app.exec();
        MainLoop::Kill();
    }
    Py_Finalize();
    return ret;
}


