#include <QMainWindow>
#include <boost/python.hpp>

class LimeWindow : public QMainWindow
{
public:
    LimeWindow(boost::python::object pyWindow);
    virtual ~LimeWindow();

private:
    boost::python::object mWindow;
};
