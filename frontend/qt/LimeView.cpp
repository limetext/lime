#include "LimeView.h"
#include <QResizeEvent>
#include <QKeyEvent>
#include <QTextLayout>
#include <QTextBlock>
#include <QPainter>
#include <QPixmap>
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
    QString font_face = QString::fromLatin1(extract<const char*>(boost::python::str(settings("font_face"))));
    float font_size = extract<float>(settings("font_size"));
    font_face = font_face.replace(" Regular", "");

    doc.setDefaultFont(QFont(font_face, font_size));
    doc.setHtml(str);
    setMinimumSize(QSize(doc.size().width(), doc.size().height()));
    timer.start();
    setFocusPolicy(Qt::ClickFocus);

    mPixmap = new QPixmap(doc.size().width(), doc.size().height());
    object bg = lime.attr("background_color")();
    QPainter p(mPixmap);
    p.fillRect(0, 0, doc.size().width(), doc.size().height(), QColor::fromRgb(extract<int>(bg[0]), extract<int>(bg[1]), extract<int>(bg[2])));
    doc.drawContents(&p);
}
LimeView::~LimeView()
{
    delete mPixmap;
}

void LimeView::resizeEvent(QResizeEvent* event)
{
    int width = event->size().width();
//    if (width != (int) doc.textWidth())
    {
//        doc.setTextWidth(width);
        setMinimumSize(QSize(doc.size().width(), doc.size().height()));
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
#if 1
        doc.drawContents(&painter, ev->rect());
#else
        QTextBlock block = doc.begin();
        int j = 0;
        int vis = 0;
        while (block.isValid())
        {
            const QTextLayout* layout = block.layout();
            QPointF pos(layout->position());
            float height = layout->lineAt(0).height();
            float off = pos.y()+height;

            if (off > ev->rect().y())
            {
                //printf("%128s - %d, %f, %d, %d\n", block.text().toAscii().constData(), j, pos.y(), ev->rect().y(), ev->rect().bottom());
                layout->draw(&painter, QPointF());
                vis += 1;
            }
            if (off > ev->rect().bottom())
                break;
            block = block.next();
            j++;
        }

        printf("%d, %d, %d, %d, %d\n", vis, ev->rect().x(), ev->rect().y(), ev->rect().width(), ev->rect().height());

//        painter.drawPixmap(0, 0, doc.size().width(), doc.size().height(), *mPixmap);
#endif
//        printf("drawing to: %d, %d, %d, %d\n", ev->rect().left(), ev->rect().top(), ev->rect().width(), ev->rect().height());

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

