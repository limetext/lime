#ifndef __INCLUDED_LIMEVIEWWIDGET_H
#define __INCLUDED_LIMEVIEWWIDGET_H
#include <QWidget>
#include <boost/python.hpp>

class LimeView;
class QScrollArea;
class LimeViewWidget : public QWidget
{
public:
    LimeViewWidget(boost::python::object o);

    LimeView* GetLimeView() const
    {
        return mLimeView;
    }
    QScrollArea* GetScrollArea() const
    {
        return mScrollArea;
    }
private:

    LimeView* mLimeView;
    QScrollArea *mScrollArea;

};
#endif
