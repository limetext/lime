#include "MainLoop.h"
#include "LimeWindow.h"

#include <boost/python.hpp>

using namespace boost::python;


static MainLoop* instance = NULL;

MainLoop* MainLoop::GetInstance()
{
    if (!instance)
        instance = new MainLoop;
    return instance;
}

void MainLoop::Kill()
{
    delete instance;
    instance = NULL;
}

MainLoop::MainLoop()
{
    try
    {
        object mainModule = import("__main__");
        object dict = mainModule.attr("__dict__");
        exec("\n\
import sys \n\
sys.path.append(\"../frontend/qt/\")\n\
import lime", dict, dict);
        mLime = dict["lime"];
        mBackend = mLime.attr("backend");
        mEditor = mLime.attr("editor");
    }
    catch (std::exception &e)
    {
        printf("Caught exception: %s\n", e.what());
    }
    catch(boost::python::error_already_set& e)
    {
        PyErr_Print();
    }
    QObject::connect(&mainTimer, SIGNAL(timeout()), this, SLOT(update()));
    mainTimer.start(100);
}

MainLoop::~MainLoop()
{
    mainTimer.stop();
}

void MainLoop::update()
{
    bool ret = extract<bool>(mEditor.attr("update")());
    if (ret)
    {
        mainTimer.setInterval(0);
    }
    else
    {
        mainTimer.setInterval(100);
    }
}

object MainLoop::GetEditor()
{
    return mEditor;
}

object MainLoop::GetLime()
{
    return mLime;
}
object MainLoop::GetBackend()
{
    return mBackend;
}

