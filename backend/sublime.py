import sys
import os.path
import platform as plat
sys.path.append("%s/3rdparty/appdirs/lib" % os.path.dirname(os.path.abspath(__file__)))
import appdirs
import backend

def packages_path():
    app_dir = appdirs.user_data_dir("Sublime Text 2", "", roaming=True)
    return "%s%sPackages" % (app_dir, os.path.sep)

def platform():
    lut = {"Darwin": "osx", "Linux": "linux", "Windows": "windows"}
    return lut[plat.system()]

def arch():
    return plat.machine()


class Region:
    def __init__(self, a, b):
        self.a = a
        self.b = b

    def __repr__(self):
        return "(%d, %d)" % (self.a, self.b)

    def intersects(self, other):
        ob = other.begin()
        oe = other.end()

        return self.contains(ob) or \
               self.contains(oe)

    def begin(self):
        return min(self.a, self.b)

    def end(self):
        return max(self.a, self.b)

    def contains(self, point):
        return point > self.start() and point <= self.end()

    def empty(self):
        return self.a == self.b

    def size(self):
        return self.end()-self.begin()


