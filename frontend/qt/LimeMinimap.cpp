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
    QObject::connect(mView->GetScrollArea()->horizontalScrollBar(), SIGNAL(valueChanged(int)), this, SLOT(repaint()));
}

void LimeMinimap::paintEvent(QPaintEvent * e)
{
    QScrollArea *sa = mView->GetScrollArea();
    QScrollBar *s = sa->verticalScrollBar();
    QScrollBar *hs = sa->horizontalScrollBar();
    LimeView * v = mView->GetLimeView();
    float prog = s->value()/(float)(s->maximum());
    float prog2 = hs->value()/(float)(hs->maximum());

    QPainter painter(this);
    float h = doc.size().height()-height();
    float w = doc.size().width()-width();
    if (h < 0)
        h = 0;
    if (w < 0)
        w = 0;
    QRect r = e->rect();
    painter.translate(0,  -prog*h);
    r.setRect(0, prog*h, width(), height());
    doc.drawContents(&painter, r);

    float visHeight = height()*float(doc.size().height())/v->doc.size().height();
    float visWidth = width();

    painter.resetTransform();
    painter.fillRect(e->rect(), QColor(0,0,0,255*0.25));

    r.setRect(prog2*visWidth, prog*(height()-visHeight), visWidth, visHeight);
    painter.fillRect(r, QColor(255,255,255,255*0.25));
}

void LimeMinimap::resizeEvent(QResizeEvent* event)
{
    QScrollArea *sa = mView->GetScrollArea();
    LimeView * v = mView->GetLimeView();

    float scale = sa->width()/(float)v->width();
    printf("scale: %f\n", scale);
    if (scale > 1)
        scale = 1;

    float width = std::max(scale*sa->width()*minimapsize, 50.0f);
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

void LimeMinimap::mousePressEvent(QMouseEvent *e)
{
    float prog = e->y()/float(height());
    QScrollArea *sa = mView->GetScrollArea();
    QScrollBar *s = sa->verticalScrollBar();
    LimeView * v = mView->GetLimeView();
    float visHeight = height()*float(doc.size().height())/v->doc.size().height();
    float prog2 = s->value()/float(s->maximum());

    float mul = prog*prog2*visHeight/height();
    printf("value: %f, %f, %f, %f\n", prog, prog2, visHeight/height(), mul);
    s->setValue(mul*s->maximum());
}
void LimeMinimap::mouseMoveEvent(QMouseEvent *e)
{

}
