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
    QString font_face = QString::fromAscii(extract<const char*>(boost::python::str(settings("font_face"))));
    int font_size = extract<int>(settings("font_size"));
    font_face = font_face.replace(" Regular", "");

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
//    if (width != (int) doc.textWidth())
    {
//        doc.setTextWidth(width);
        setMinimumSize(QSize(0, doc.size().height()));
    }
}

void LimeView::keyPressEvent(QKeyEvent* ke)
{
    if (ke->key() > 160000)
    {
        // TODO: we don't want to process presses that are just modifiers
        return;
    }
    Qt::KeyboardModifiers mod = ke->modifiers();
    MainLoop *loop = MainLoop::GetInstance();
    object editor = loop->GetEditor();
    try
    {
        printf("ke->text().unicode()[0].unicode(): %d\n", ke->key() & Qt::MODIFIER_MASK);
        Q_ASSERT(ke->text().length() == 1);
        object backend = loop->GetBackend();
        object keypress = backend.attr("KeyPress")(
                QChar(ke->key()).toLower().unicode(),
                ke->text().unicode()[0].unicode(),
#if __APPLE__
                mod.testFlag(Qt::ControlModifier),
                mod.testFlag(Qt::MetaModifier),
                mod.testFlag(Qt::ShiftModifier),
                mod.testFlag(Qt::AltModifier));
#else
                mod.testFlag(Qt::MetaModifier),
                mod.testFlag(Qt::ControlModifier),
                mod.testFlag(Qt::ShiftModifier),
                mod.testFlag(Qt::AltModifier));
#endif
        editor.attr("keyEvent")(keypress);
    }
    catch (...)
    {
        PyErr_Print();
    }
}
void LimeView::paintEvent(QPaintEvent* ev)
{
    static int count = 0;
    //if (++count < 2)
    {
        QPainter painter(this);
        doc.drawContents(&painter, ev->rect());
//        printf("drawing to: %d, %d, %d, %d\n", ev->rect().left(), ev->rect().top(), ev->rect().width(), ev->rect().height());
        /*
        QFont f = doc.defaultFont();
        qreal old = f.pointSizeF();
        float minimapsize = 0.1;
        f.setPointSizeF(old*minimapsize);
        doc.setDefaultFont(f);
        QRect r = ev->rect();
        painter.translate(QPoint(r.right()-r.width()*minimapsize, 0));
        printf("%d, %d, %d, %d\n", r.left(), r.top(), r.width(), r.height());
        doc.drawContents(&painter, ev->rect());
        f.setPointSizeF(old);
        doc.setDefaultFont(f);
        */

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

