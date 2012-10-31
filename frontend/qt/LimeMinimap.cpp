#include "LimeMinimap.h"
#include "LimeView.h"
#include "LimeViewWidget.h"
#include <QPaintEvent>
#include <QResizeEvent>
#include <QPainter>
#include <QScrollArea>
#include <QScrollBar>
#include <QRegExp>
#include <QWheelEvent>
#include <QApplication>

static const float minimapsize = 0.25;
LimeMinimap::LimeMinimap(LimeViewWidget* v)
: mView(v)
{
    QTextDocument &origDoc = mView->GetLimeView()->doc;
    QFont f = origDoc.defaultFont();
    f.setPointSizeF(f.pointSizeF()*minimapsize);
    doc.setDefaultFont(f);
    QString html = origDoc.toHtml().replace(QRegExp("font-\\w+:[^;]+;"), "");
    doc.setHtml(html);
    QObject::connect(mView->GetScrollArea()->verticalScrollBar(), SIGNAL(valueChanged(int)), this, SLOT(repaint()));
}

void LimeMinimap::paintEvent(QPaintEvent * e)
{
    QScrollArea *sa = mView->GetScrollArea();
    QScrollBar *s = sa->verticalScrollBar();
    float prog = s->value()/(float)(s->maximum());

    // if (mView->GetLimeView()->doc.textWidth()*minimapsize != (doc.textWidth()))
    // {
    //     doc.setTextWidth(mView->GetLimeView()->doc.textWidth()*minimapsize);
    //     printf("%f, %f, %f\n", mView->GetLimeView()->doc.textWidth(), doc.textWidth(), minimapsize*mView->GetLimeView()->doc.textWidth());
    // }
    QFont f = doc.defaultFont();
    QPainter painter(this);
    float h = doc.size().height()-height();
    QRect r = e->rect();
    if (h > 0)
    {
        painter.translate(0,  -prog*h);
        r.setRect(0, prog*h, width(), height());
    }
    doc.drawContents(&painter, r);

    float visHeight = height()*float(doc.size().height())/mView->GetLimeView()->doc.size().height();

    painter.resetTransform();
    painter.fillRect(e->rect(), QColor(0,0,0,255*0.25));

    r.setRect(0, prog*(height()-visHeight), width(), visHeight);
    painter.fillRect(r, QColor(255,255,255,255*0.25));
}

void LimeMinimap::resizeEvent(QResizeEvent* event)
{
    LimeView * v = mView->GetLimeView();
    float width = std::max(v->width()*minimapsize, 100.0f);
    width = std::min(150.0f, width);
    setMinimumSize(width, 0);

}

void LimeMinimap::wheelEvent(QWheelEvent* event)
{
    QPoint pos(10, 10);
    QWheelEvent e(pos, mView->GetScrollArea()->mapToGlobal(pos), event->delta(), event->buttons(), event->modifiers());
    printf("sending event: %d, %d, %d\n", e.pos().x(), e.pos().y(), e.delta());
    QApplication::sendEvent(mView->GetScrollArea(), &e);
}
