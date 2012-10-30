#include "LimeWindow.h"
#include <QStatusBar>
#include "LimeView.h"
#include <QScrollArea>
#include "MainLoop.h"

using namespace boost::python;

static void view_added(object view)
{
    object window = view.attr("window")();
    LimeWindow* qtLimeWindow = extract<LimeWindow*>(window.attr("qtLimeWindow"));
    QTabWidget *qt = static_cast<QTabWidget*>(qtLimeWindow->centralWidget());
    char * str = extract<char *>(view.attr("file_name")());
    LimeView *v = new LimeView(view);
    QScrollArea *s = new QScrollArea();
    s->setWidgetResizable(true);
    s->setWidget(v);
    qt->addTab(s, str);
}

LimeWindow::LimeWindow(object pyWindow) : QMainWindow(), mWindow(pyWindow)
{
    class_<LimeWindow, boost::noncopyable>("LimeWindow", no_init);
    QString style;
    object lime = MainLoop::GetInstance()->GetLime();
    style = extract<const char *>(lime.attr("cssify_theme")());
    setStyleSheet(style);

    QStatusBar* bar = new QStatusBar;
    bar->showMessage("This is a status bar");
    setStatusBar(bar);
    QTabWidget* qt = new QTabWidget;
    qt->setAutoFillBackground(true);
    qt->setTabsClosable(true);
    qt->setMovable(true);

    printf("final style is:\n%s\n", style.toAscii().constData());

    setCentralWidget(qt);
    mWindow.attr("view_added_event") += view_added;
}

LimeWindow::~LimeWindow()
{

}
