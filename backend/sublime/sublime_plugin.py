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

class Command(object):
    def is_enabled(self, args=None):
        return True

    def is_visible(self, args=None):
        return True

class ApplicationCommand(Command):
    pass

class WindowCommand(Command):
    def __init__(self, wnd):
        self.window = wnd

class TextCommand(Command):
    def __init__(self, view):
        self.view = view

    def run_(self, edit_token, kwargs):
        if kwargs and 'event' in kwargs:
            del kwargs['event']

        if kwargs:
            self.run(edit_token, **kwargs)
        else:
            self.run(edit_token)

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
    def cmdname(name):
        if name.endswith("Command"):
            name = name[:-7]
        ret = ""
        for c in name:
            l = c.lower()
            if c != l and len(ret) > 0:
                ret += "_"
            ret += l
        return ret
    print "Loading plugin %s" % module
    try:
        module = importlib.import_module(module)
        for item in inspect.getmembers(module):
            if type(EventListener) != type(item[1]):
                continue

            try:
                cmd = cmdname(item[0])
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
                    sublime.register(cmd, sublime.TextCommandGlue(item[1]))
                elif issubclass(item[1], WindowCommand):
                    sublime.register(cmd, sublime.WindowCommandGlue(item[1]))
                elif issubclass(item[1], ApplicationCommand):
                    sublime.register(cmd, sublime.ApplicationCommandGlue(item[1]))
            except:
                print "Skipping registering %s: %s" % (cmd, sys.exc_info()[1])
        if "plugin_loaded" in dir(module):
            module.plugin_loaded()
    except:
        traceback.print_exc()


class MyLogger:
    def write(self, data):
        v = sublime.console()
        try:
            e = v.begin_edit()
            v.insert(e, v.size(), data)
        finally:
            v.end_edit(e)

sys.stdout = sys.stderr = MyLogger()
