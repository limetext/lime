import os
import os.path
import inspect
import traceback
import sublime

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

def reload_plugin(filename):
    print "Loading plugin %s" % filename
    oldpath = os.getcwd()
    path = os.path.dirname(os.path.abspath(filename))
    try:
        os.chdir(path)
        filename = os.path.relpath(filename, path)
        module = os.path.splitext(filename)[0]
        module = __import__(module)
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
    finally:
        os.chdir(oldpath)


