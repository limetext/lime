import sys
sys.path.append("3rdparty/appdirs/lib")
import appdirs
import os.path

def packages_path():
    app_dir = appdirs.user_data_dir("Sublime Text 2", "", roaming=True)
    return "%s%sPackages" % (app_dir, os.path.sep)

print packages_path()
