#include <QObject>
#include <QTimer>

#ifndef Q_MOC_RUN
#include <boost/python.hpp>
using namespace boost::python;
#endif

class MainLoop : public QObject
{
    Q_OBJECT
public:
    static MainLoop* GetInstance();
    static void Kill();

    object GetEditor();
    object GetLime();
    object GetBackend();

public slots:
    void update();

private:
    MainLoop();
    virtual ~MainLoop();
    object mEditor;
    object mLime;
    object mBackend;
    QTimer mainTimer;
};
