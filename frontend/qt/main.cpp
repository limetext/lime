#include <QApplication>
#include <QTextDocument>
#include <QWidget>
#include <QPainter>
#include <QScrollArea>
#include <QFile>
#include <QPaintEvent>
#include <QTextCursor>
#include <QTextBlock>
#include <QTextLayout>
#include <QTextLine>
#include <QElapsedTimer>
#include <QTimer>
#include <QTabWidget>
#include <QTextEdit>
#include <QVBoxLayout>
#include <QStatusBar>
#include <QMainWindow>
#include <QStaticText>
#include <Python.h>
#include <math.h>
#include "mainloop.h"

static const char *initdata = "\n\
import sys\n\
sys.path.append(\"../frontend/qt/\")\n\
import lime\n\
";

PyObject *dict = NULL;
void MainLoop::update()
{
    PyRun_String("editor.update()", Py_single_input, dict, dict);
}

int main(int argc, char** argv)
{
    QApplication app(argc, argv);
    Py_Initialize();

    PyObject *dict = PyDict_New();
    PyObject *obj = PyRun_String(initdata, Py_single_input, dict, dict);

    class LimeViewWidget : public QWidget
    {
    public:
        LimeViewWidget(QWidget* parent = 0) : QWidget(parent), pixmap(0)
        {
            QFile f("/Users/quarnster/code/lime/frontend/data.txt");
            f.open(QIODevice::ReadOnly);
            char* data = new char[f.size() + 1];
            f.read(data, f.size());
            data[f.size()] = '\0';
            QString data2(data);
            data2 = data2.replace("\n\n", "<p>");
            QString html;
            html += data2;
            doc.setHtml(html);
            st.setText(html);
            st.setTextFormat(Qt::RichText);
            f.close();
            delete[] data;
            doc.setDefaultFont(QFont("Menlo", 11));
            const QSizeF& q = doc.size();
            resize(600, q.height());
            const QSize& q2 = size();

            printf("%d, %d\n", q2.width(), q2.height());
            printf("here\n");

            for (int i = 0; i < 3; i++)
            {
                cursors[i] = new QTextCursor(&doc);
                cursors[i]->setPosition((i + 1) * 500);
                cursors[i]->select(QTextCursor::WordUnderCursor);
                printf("%s\n", cursors[i]->selectedText().toLatin1().constData());
                cursors[i]->setPosition((i + 1) * 500);
            }

            timer.start();
            setFocusPolicy(Qt::ClickFocus);
        }
        virtual ~LimeViewWidget()
        {
            delete pixmap;
        }
        QTextCursor* cursors[3];
        QElapsedTimer timer;
        QPixmap* pixmap;
        QTextDocument doc;
        QStaticText st;

        virtual void resizeEvent(QResizeEvent* event)
        {
            int width = event->size().width();
            if (width != (int) doc.textWidth())
            {
                doc.setDefaultFont(QFont("Menlo", 11));
                doc.setTextWidth(width);
                delete pixmap;
                setMinimumSize(QSize(0, doc.size().height()));
            }
        }

        virtual void keyPressEvent(QKeyEvent* ke)
        {
            //window()->setWindowState(window()->windowState() ^ Qt::WindowFullScreen);
            Qt::KeyboardModifiers mod = ke->modifiers();
            printf("mod: %d%d%d%d, key: %d (%c), scancode: %d (%c), text: %s\n", mod.testFlag(Qt::MetaModifier), mod.testFlag(Qt::ControlModifier), mod.testFlag(Qt::AltModifier), mod.testFlag(Qt::ShiftModifier), ke->key(), ke->key(), ke->nativeVirtualKey(), ke->nativeVirtualKey(), ke->text().constData());
        }
        virtual void paintEvent(QPaintEvent* ev)
        {
            static int count = 0;
            //if (++count < 2)
            {
                QPainter painter(this);
                doc.drawContents(&painter, ev->rect());
                float brightness = 128 + 127 * sin(timer.elapsed() * 1.0 / (150.0));
                painter.setPen(QColor::fromRgb(0xff, 0xff, 0xff, brightness));

                for (int i = 0; i < 3; i++)
                {
                    int pos = cursors[i]->position();

                    const QTextBlock& block = doc.findBlock(pos);
                    const QTextLayout* layout = block.layout();
                    pos -= block.position();
                    layout->drawCursor(&painter, QPointF(0, 0), pos, 2);
                }
            }
        }

    };
    QTextEdit* t = new QTextEdit;
    t->setText("Hello World!");

    QScrollArea* s = new QScrollArea;
    LimeViewWidget* w = new LimeViewWidget;
    s->resize(600, 400);
    s->setWidget(w);
    s->setWidgetResizable(true);
    QTabWidget* qt = new QTabWidget;
    qt->addTab(t, "Hello");
    qt->addTab(s, "World");
    qt->setTabsClosable(true);
    //    qt->setDocumentMode(true);
    QString style;
    style += "QTabWidget {  position: absolute; left: 0; right: 0; background-color: yellow; padding-right: 0; }";
    style += "QTabWidget::pane {  background-color: blue;  }";
    style += "QTabWidget::tab-bar {subcontrol-origin: border;  position: absolute; left: 0; right: 0; background-color: blue; padding-right: 0; }";
    style += "QTabBar {  position: absolute; left: 0; right: 0;  background: red; }";
    style += "LimeViewWidget { text: \"Some text\"; renderType: Text.NativeRendering; }";

    /*
        style += "QTabBar::tear {  background: blue; margin: 0px; border: 2px; border-style:solid; border-image: url(../../3rdparty/Theme - Soda/Soda Dark/tabset-background.png) 2 2 2 2 stretch stretch; right: 0; }";
        style += "QTabWidget::tab-bar {  position: absolute; subcontrol-position: left; background-color: red; left: 0; right: 0; margin: 0px; border: 1px;}";
        style += "QTabBar::tab { border-image: url(../../3rdparty/Theme - Soda/Soda Dark/tab-inactive.png); border: 4px; }";
        style += "QTabBar::tab:hover { border-image: url(../../3rdparty/Theme - Soda/Soda Dark/tab-hover.png); }";
        style += "QTabBar::tab:selected { border-image: url(../../3rdparty/Theme - Soda/Soda Dark/tab-active.png); }";
        style += "QTabBar::close-button { image: url(../../3rdparty/Theme - Soda/Soda Dark/tab-close.png); }";
        style += "QTabBar::close-button:selected { image: url(../../3rdparty/Theme - Soda/Soda Dark/tab-close-inactive.png); }";
    */
    qt->setStyleSheet(style);

    QStatusBar* bar = new QStatusBar;
    bar->showMessage("This is a status bar");

    QMainWindow mw;
    mw.resize(600, 400);
    mw.setCentralWidget(qt);
    mw.setStatusBar(bar);
    mw.show();

    QTimer mainTimer;
    MainLoop d;
    QObject::connect(&mainTimer, SIGNAL(timeout()), &d, SLOT(update()));
    mainTimer.start(0);

    int ret = app.exec();
    PyRun_String("editor.exit(0)", Py_single_input, dict, dict);
    //Py_DECREF(obj);
    Py_DECREF(dict);
    Py_Finalize();
    return ret;
}


