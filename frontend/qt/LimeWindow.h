#include <QMainWindow>

#ifndef Q_MOC_RUN
#include <boost/python.hpp>
#endif

class LimeTabWidget : public QTabWidget
{
    Q_OBJECT
public:
    LimeTabWidget();

protected:
    virtual void resizeEvent(QResizeEvent * e);
    virtual void paintEvent(QPaintEvent * e);
private slots:
    void fixSize();
};

class LimeWindow : public QMainWindow
{
public:
    LimeWindow(boost::python::object pyWindow);
    virtual ~LimeWindow();

private:
    boost::python::object mWindow;
};
