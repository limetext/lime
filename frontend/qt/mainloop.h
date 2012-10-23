#include <QObject>

class MainLoop : public QObject
{
    Q_OBJECT
public slots:
    void update();
};
