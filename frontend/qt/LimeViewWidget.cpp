#include "LimeViewWidget.h"
#include "LimeView.h"
#include "LimeMinimap.h"
#include <QScrollArea>
#include <QHBoxLayout>
#include <QWheelEvent>
#include <QApplication>
#include <QPoint>

LimeViewWidget::LimeViewWidget(boost::python::object view)
{
    mLimeView = new LimeView(view);
    mScrollArea = new QScrollArea(this);
    mScrollArea->setVerticalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    mScrollArea->setHorizontalScrollBarPolicy(Qt::ScrollBarAlwaysOff);
    mScrollArea->setWidgetResizable(true);
    mScrollArea->setWidget(mLimeView);
    QHBoxLayout *layout = new QHBoxLayout;
    layout->setSpacing(0);
    layout->setMargin(0);
    layout->addWidget(mScrollArea);
    layout->addWidget(new LimeMinimap(this));
    setLayout(layout);
}
