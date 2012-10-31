#include "LimeWindow.h"
#include <QStatusBar>
#include "LimeViewWidget.h"
#include <QTabBar>
#include <QResizeEvent>
#include <QPaintEvent>
#include "MainLoop.h"

using namespace boost::python;

class Test : public QWidget
{
public:
    Test()
    {
        resize(1000, 20);
        setAutoFillBackground(true);
    }
};

LimeTabWidget::LimeTabWidget()
{
    setAutoFillBackground(true);
    setTabsClosable(true);
    setMovable(true);
    setAcceptDrops(true);
    setDocumentMode(true);
    //tabBar()->setDrawBase(true);
    //tabBar()->setAutoFillBackground(true);
    QObject::connect(this, SIGNAL(currentChanged(int)), this, SLOT(fixSize()));
}

void LimeTabWidget::resizeEvent(QResizeEvent * e)
{
    QTabWidget::resizeEvent(e);
    fixSize();
}
void LimeTabWidget::paintEvent(QPaintEvent * e)
{
    QTabWidget::paintEvent(e);
    fixSize();
}

void LimeTabWidget::fixSize()
{
//    tabBar()->resize(size().width(), tabBar()->size().height());
}

static void view_added(object view)
{
    object window = view.attr("window")();
    LimeWindow* qtLimeWindow = extract<LimeWindow*>(window.attr("qtLimeWindow"));
    QTabWidget *qt = static_cast<QTabWidget*>(qtLimeWindow->centralWidget());
    char * str = extract<char *>(view.attr("file_name")());
    LimeViewWidget *v = new LimeViewWidget(view);
    qt->addTab(v, str);
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
    QTabWidget* qt = new LimeTabWidget;
    qt->setCornerWidget(new Test(), Qt::TopRightCorner);


    printf("final style is:\n%s\n", style.toAscii().constData());

    setCentralWidget(qt);
    mWindow.attr("view_added_event") += view_added;
}

LimeWindow::~LimeWindow()
{

}
