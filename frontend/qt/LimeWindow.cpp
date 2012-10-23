#include "LimeWindow.h"
#include <QStatusBar>
#include "LimeView.h"
#include <QScrollArea>

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
    QStatusBar* bar = new QStatusBar;
    bar->showMessage("This is a status bar");
    setStatusBar(bar);
    QTabWidget* qt = new QTabWidget;
    qt->setTabsClosable(true);
    QString style;
    style += "QTabWidget {  position: absolute; left: 0; right: 0; background-color: yellow; padding-right: 0; }";
    style += "QTabWidget::pane {  background-color: blue;  }";
    style += "QTabWidget::tab-bar {subcontrol-origin: border;  position: absolute; left: 0; right: 0; background-color: blue; padding-right: 0; }";
    style += "QTabBar {  position: absolute; left: 0; right: 0;  background: red; }";
    style += "LimeViewWidget { text: \"Some text\"; renderType: Text.NativeRendering; }";

    qt->setStyleSheet(style);

    setCentralWidget(qt);
    mWindow.attr("view_added_event") += view_added;
}

LimeWindow::~LimeWindow()
{

}
