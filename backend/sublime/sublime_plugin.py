import os
import os.path
import inspect
import traceback
import imp
import sublime
import sys
import importlib

class __Event:
    def __init__(self):
        self.__observers = []

    def __call__(self, *args):
        for observer in self.__observers:
            try:
                observer(*args)
            except:
                traceback.print_exc()

    def __iadd__(self, observer):
        self.__observers.append(observer)
        return self

    def __isub__(self, observer):
        self.__observers.remove(observer)
        return self


on_load = __Event()
on_new = __Event()

application_commands = {}
window_commands = {}
text_commands = {}

class Command(object):
    pass

class ApplicationCommand(Command):
    pass

class WindowCommand(Command):
    pass

class TextCommand(Command):
    pass

class EventListener(object):
    pass


def fn(fullname):
    paths = fullname.split(".")
    paths = "/".join(paths)
    for p in sys.path:
        f = os.path.join(p, paths)
        if os.path.exists(f):
            return f
        f += ".py"
        if os.path.exists(f):
            return f
    return None

class __myfinder:
    class myloader(object):
        def load_module(self, fullname):
            if fullname in sys.modules:
                return sys.modules[fullname]
            f = fn(fullname)
            if not f.endswith(".py"):
                m = imp.new_module(fullname)
                m.__path__ = f
                sys.modules[fullname] = m
                return m
            return imp.load_source(fullname, f)

    def find_module(self, fullname, path=None):
        f = fn(fullname)
        if f != None:
            return self.myloader()



sys.meta_path = [__myfinder()]

def reload_plugin(module):
    print "Loading plugin %s" % module
    try:
        module = importlib.import_module(module)
        for item in inspect.getmembers(module):
            if type(EventListener) != type(item[1]):
                continue

            try:
                if issubclass(item[1], EventListener):
                    def add(inst, listname):
                        toadd = getattr(inst, listname, None)
                        if toadd:
                            l = eval(listname)
                            l += toadd
                    inst = item[1]()
                    add(inst, "on_load")
                    add(inst, "on_new")
                elif issubclass(item[1], TextCommand):
                    text_commands[item[0]] = item[1]
                elif issubclass(item[1], WindowCommand):
                    window_commands[item[0]] = item[1]
                elif issubclass(item[1], ApplicationCommand):
                    application_commands[item[0]] = item[1]
            except:
                traceback.print_exc()
    except:
        traceback.print_exc()
