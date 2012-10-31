#include <QWidget>
#include <QElapsedTimer>
#include <QTextDocument>
#include <boost/python.hpp>

class LimeView : public QWidget
{
public:
    LimeView(boost::python::object view);
    virtual ~LimeView();

    virtual void resizeEvent(QResizeEvent* event);
    virtual void keyPressEvent(QKeyEvent* ke);
    virtual void paintEvent(QPaintEvent* ev);
private:
    boost::python::object mView;

    QElapsedTimer timer;
    QTextDocument doc;

    friend class LimeMinimap;
};
