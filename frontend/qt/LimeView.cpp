#include "LimeView.h"
#include <QResizeEvent>
#include <QKeyEvent>
#include <QTextLayout>
#include <QTextBlock>
#include <QPainter>
#include <math.h>
#include "MainLoop.h"


LimeView::LimeView(boost::python::object view) : QWidget(), mView(view)
{
    object lime(MainLoop::GetInstance()->GetLime());
    object style = lime.attr("create_stylesheet")();
    object data = lime.attr("transform_scopes")(view);
    data = style + data;
    const char * str = extract<const char *>(data);

    object settings = view.attr("settings")().attr("get");
    const char *font_face = extract<const char*>(boost::python::str(settings("font_face")));
    int font_size = extract<int>(settings("font_size"));

    doc.setDefaultFont(QFont(font_face, font_size));
    doc.setHtml(str);
    setMinimumSize(QSize(0, doc.size().height()));
    timer.start();
    setFocusPolicy(Qt::ClickFocus);
}
LimeView::~LimeView()
{
}

void LimeView::resizeEvent(QResizeEvent* event)
{
    int width = event->size().width();
    if (width != (int) doc.textWidth())
    {
        doc.setTextWidth(width);
        setMinimumSize(QSize(0, doc.size().height()));
    }
}

void LimeView::keyPressEvent(QKeyEvent* ke)
{
    //window()->setWindowState(window()->windowState() ^ Qt::WindowFullScreen);
    Qt::KeyboardModifiers mod = ke->modifiers();
    printf("mod: %d%d%d%d, key: %d (%c), scancode: %d (%c), text: %s\n", mod.testFlag(Qt::MetaModifier), mod.testFlag(Qt::ControlModifier), mod.testFlag(Qt::AltModifier), mod.testFlag(Qt::ShiftModifier), ke->key(), ke->key(), ke->nativeVirtualKey(), ke->nativeVirtualKey(), ke->text().constData());
}
void LimeView::paintEvent(QPaintEvent* ev)
{
    static int count = 0;
    //if (++count < 2)
    {
        QPainter painter(this);
        doc.drawContents(&painter, ev->rect());
        float brightness = 128 + 127 * sin(timer.elapsed() * 1.0 / (150.0));
        painter.setPen(QColor::fromRgb(0xff, 0xff, 0xff, brightness));

        /*
        for (int i = 0; i < 3; i++)
        {
            int pos = cursors[i]->position();

            const QTextBlock& block = doc.findBlock(pos);
            const QTextLayout* layout = block.layout();
            pos -= block.position();
            layout->drawCursor(&painter, QPointF(0, 0), pos, 2);
        }
        */
    }
}

