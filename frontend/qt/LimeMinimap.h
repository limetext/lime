#ifndef __INCLUDED_LIME_MINIMAP_H
#define __INCLUDED_LIME_MINIMAP_H

#include <QWidget>
#include <QTextDocument>

class LimeViewWidget;
class LimeMinimap : public QWidget
{
public:
    LimeMinimap(LimeViewWidget*);
protected:
    virtual void paintEvent(QPaintEvent * e);
    virtual void resizeEvent(QResizeEvent * e);
    virtual void wheelEvent(QWheelEvent *e);
private:
    LimeViewWidget* mView;
    QTextDocument doc;
};

#endif
