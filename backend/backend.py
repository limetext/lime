import re
import plistlib
import traceback
import sys
import json
import os
import os.path
sys.path.append("%s/3rdparty/appdirs/lib" % os.path.dirname(os.path.abspath(__file__)))
sys.path.append(".")
import appdirs
import pickle
import gzip
import time
import Queue

def loadjson( name):
    f = open(name)
    data = f.read()
    f.close()
    regex1 = re.compile(r"(//[^\n]*)(?=\n)", re.MULTILINE|re.DOTALL)
    regex2 = re.compile(r"(?<!\\)(\".*?(?<!\\)\")", re.MULTILINE|re.DOTALL)

    off = 0
    while True:
        m1 = regex1.search(data, off)
        m2 = regex2.search(data, off)
        if m1 and (not m2 or (m1.start() < m2.start())):
            data = "%s%s" % (data[:m1.start()], data[m1.end():])
            off = m1.start()
        elif m2:
            off = m2.end()
        if not m1:
            break
    regex = re.compile(r"/\*.*?\*/", re.MULTILINE|re.DOTALL)
    data = regex.sub("", data)
    try:
        data = json.loads(data)
    except:
        traceback.print_exc()
        print data
        data = None
    return data

def verify(scopes):
    if len(scopes) > 1:
        a = scopes[0].region
        for b in scopes[1:]:
            #print a, b.region
            b = b.region
            assert a.end() == b.begin()
            a = b

def singleton(cls):
    instances = {}
    def getinstance():
        if cls not in instances:
            instances[cls] = cls()
        return instances[cls]
    return getinstance


@singleton
class Editor:

    class __Event:
        def __init__(self):
            self.__observers = []

        def __call__(self, *args):
            print "event triggered"
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


    class __Window:
        def __init__(self):
            self.__views = {}
            e = Editor()
            self.__settings = e._Editor__Settings(e._Editor__settings)
            # TODO: project settings

        def open_file(self, name, flags=0):
            # TOOD: handle sublime.ENCODED_POSITION
            if name not in self.__views:
                e = Editor()
                self.__views[name] = e._Editor__View(self, name)
            return self.__views[name]

        def views(self):
            return self.__views.values()

    class __View:
        class __Edit:
            pass

        def __init__(self, window, name=None):
            self.__window = window
            self.__scratch = False
            try:
                e = Editor()
                self.__settings = e._Editor__Settings(window._Window__settings)
                self.__syntax = None
                self.__settings.add_on_change("__lime_View", self.__settings_changed)
                self.__set_file(name)
            except:
                traceback.print_exc()

        def begin_edit(self):
            return self.__Edit()

        def end_edit(self, edit):
            assert edit

        def replace(self, edit, region, data):
            assert edit
            # Todo: don't commit changes until end_edit is called
            self.__buffer = "%s%s%s" % (self.__buffer[:region.begin()], data, self.__buffer[region.end():])

        def __set_file(self, name):
            self.__file = name
            if name and os.path.isfile(name):
                f = open(name)
                self.__buffer = f.read()
                f.close()
                if name:
                    self.__settings.set("syntax", Editor().get_syntax_file_for_filename(name))
            else:
                self.__buffer = ""
            sublime_plugin.on_load(self)

        def has_non_empty_selection_region(self):
            return False # Todo

        def set_status(self, key, value):
            pass # TODO

        def is_scratch(self):
            return self.__scratch

        def set_scratch(self, value):
            self.__scratch = value

        def file_name(self):
            return self.__file

        def __settings_changed(self):
            syntax = self.__settings.get("syntax")
            if self.__syntax == None or self.__syntax.file != syntax:
                e = Editor()
                self.__syntax = e.get_syntax(syntax)
                self.__scopes = self.__syntax.extract_scopes(self.__buffer)

        def __find_scope(self, point):
            for scope in self.__scopes:
                if scope.region.contains(point):
                    return scope
            return None

        def extract_scope(self, point):
            scope = self.__find_scope(point)
            return None if not scope else scope.region

        def scope_name(self, point):
            scope = self.__find_scope(point)
            return scope.name if scope else self.__syntax.get_default_scope()

        def window(self):
            return self.__window

        def settings(self):
            return self.__settings

        def size(self):
            return len(self.__buffer)

        def substr(self, region):
            return self.__buffer[region.begin():region.end()]

        def run_command(self, command, *args):
            print "wants to run command %s, %s" % (command, args)

    class __Settings:
        def __init__(self, other=None):
            self.__values = {}
            self.__on_change = {}
            if other:
                self.__update(other)

        def add_on_change(self, key, callback):
            if key not in self.__on_change:
                self.__on_change[key] = Editor()._Editor__Event()
            self.__on_change[key] += callback

        def clear_on_change(self, key):
            if key in self.__on_change:
                del self.__on_change[key]

        def get(self, name, default=None):
            if name in self.__values:
                return self.__values[name]
            return default

        def set(self, name, value):
            self.__values[name] = value
            self.__trigger_changes()

        def erase(self, name):
            if name in self.__values:
                del self.__values[name]

        def __update(self, other):
            if isinstance(other, str):
                if os.path.isfile(other):
                    self.__values.update(loadjson(other))
            else:
                self.__values.update(other.__values)
            self.__trigger_changes()

        def __trigger_changes(self):
            for key in self.__on_change:
                try:
                    self.__on_change[key]()
                except:
                    traceback.print_exc()

    class __Scope:
        def __init__(self, name, region):
            self.name   = name
            self.region = region
            assert region.a <= region.b

        def __repr__(self):
            return "%s %s" % (self.name, self.region)

    class __SyntaxPattern:
        class HackRegex:
            class HackMatchObject:
                def __init__(self, match, idx):
                    self.__groups = list(match.groups())
                    self.__groups.insert(0, match.group(0))
                    self.__start = [match.start(i)+idx for i in range(len(self.__groups))]
                    self.__end = [match.end(i)+idx for i in range(len(self.__groups))]

                def start(self, i=0):
                    return self.__start[i]

                def end(self, i=0):
                    return self.__end[i]

                def group(self, i=0):
                    return self.__groups[i]

                def groups(self):
                    return self.__groups[1:]

            def __init__(self, pattern):
                # \\G means right at the start where the previous pattern ended and
                # isn't supported by python's regex module, so we'll have to hack
                # around a bit
                self.regex = re.compile(re.sub(r"\\G", "^", pattern))
                self.pattern = pattern

            def search(self, data, idx=0):
                match = self.regex.search(data[idx:])
                if match:
                    match = self.HackMatchObject(match, idx)
                return match

        def compile(self, data, var):
            sqregex = re.compile(r"(?<!\\)\[.*?(?<!\\)\]")

            def min2(a, b):
                if a == -1:
                    return b
                if b == -1:
                    return a
                return min(a, b)

            def find_nested_end(data, idx, inc="(", dec=")"):
                count = 1
                regex = re.compile(r"(?<!\\)(%s|%s)" % (re.escape(inc), re.escape(dec)))
                while idx < len(data):
                    match = regex.search(data, idx)
                    if not match:
                       break
                    match2 = sqregex.search(data, idx)
                    idx = match.start()
                    if match2 and idx > match2.start() and idx <= match2.end():
                        idx = match2.end()
                        continue
                    if match.group(1) == inc:
                        count += 1
                    else:
                        count -= 1
                        if count == 0:
                            return idx
                return -1
            def split_lb(sub):
                off = 0
                ret = []
                match = re.match(r"(\(\?<(=|!))", sub)
                assert match
                start = match.group(1)
                sub = sub[match.end():-1]

                while off < len(sub):
                    match = re.search(r"(?<!\\)\|", sub)
                    if not match:
                        break

                    off = match.start()
                    if off != -1:
                        #might have to split it
                        match = sqregex.search(sub)
                        if match and off > match.start() and off < match.end():
                            off = match.end()
                            continue
                        ret.append("%s%s)" % (start, sub[:off]))
                        sub = sub[off+1:]
                        off = 0
                    else:
                        break
                if len(sub):
                    ret.append("%s%s)" % (start, sub))
                return ret

            try:
                data = data[var] if var in data else None
                if data:
                    # fixed regex settings
                    regex = re.compile(r"(\(\?[iLmsux]+):")
                    match = regex.search(data)
                    while match:
                        data = "%s(?:%s)%s" % (data[:match.start()], match.group(1), data[match.end():])
                        match = regex.search(data)

                    # fix lookback patterns
                    regex = re.compile(r"(\(\?<(=|!))")
                    match = regex.search(data)
                    while match:
                        start = match.start()
                        end = find_nested_end(data, start+1)
                        sub = data[start:end+1]
                        split = split_lb(sub)
                        new = "(?:%s)" % ("|".join(split))
                        data = "%s%s%s" % (data[:start], new, data[end+1:])
                        start = start + len(new)
                        match = regex.search(data, start)

                    # fix lookahead patterns
                    data = re.sub(r"(?<!\\)\(\?\>", "(?=", data)

                    # Fix named patterns
                    regex = re.compile(r"(\(\?\<(\w+)\>)")
                    data = regex.sub("(?P<\\2>", data)
                    data = re.sub(r"\\x\{([0-9a-zA-Z]{1,2})\}", r"\\x\1", data)

                    # fix multiple repeats
                    regex = re.compile(r"([+*]{2,}|\?[+*?])")
                    match = regex.search(data)
                    pos = 0
                    while match:
                        match2 = sqregex.search(data, pos)
                        if match2 and match.start() > match2.start() and match.start() < match2.end():
                            pos = match2.end()
                        else:
                            data = "%s%s%s" % (data[:match.start()], match.group(1)[-1], data[match.end():])
                        match = regex.search(data, pos)

                    if "\\G" in data:
                        data = self.HackRegex(data)
                    else:
                        data = re.compile(data, re.MULTILINE)
                setattr(self, var, data)
            except:
                setattr(self, var, None)
                print "Failed to compile regex: \"%s\" - %s" % (data, sys.exc_info()[1])

        def __init__(self, data, syntax):
            self.compile(data, "match")
            self.compile(data, "begin")
            self.compile(data, "end")

            self.captures = data["captures"] if "captures" in data else None
            self.beginCaptures = data["beginCaptures"] if "beginCaptures" in data else None
            self.endCaptures = data["endCaptures"] if "endCaptures" in data else None
            self.patterns = []
            self.syntax = syntax
            self.__cachedData = None
            self.__cachedMatch = None
            self.__cachedPat = None
            self.__cachedPatterns = None
            self.include = None
            self.hits = 0
            self.misses = 0
            if self.begin and not self.match and self.captures and not self.beginCaptures:
                self.beginCaptures = self.captures
                self.captures = None

            if "patterns" in data:
                e = Editor()
                for pattern in data["patterns"]:
                    self.patterns.append(e._Editor__SyntaxPattern(pattern, syntax))
            if "include" in data:
                self.include = data["include"]

            self.name = data["name"] if "name" in data else None
            if self.name == None and "contentName" in data:
                self.name = data["contentName"]

        def cache(self, data, pos):
            pat = None
            ret = None
            if self.__cachedData == data:
                if not self.__cachedMatch:
                    return None, None
                if self.__cachedMatch.start() > pos:
                    self.hits += 1
                    return self.__cachedPat, self.__cachedMatch
            else:
                self.__cachedPatterns = None
            if self.__cachedPatterns == None:
                self.__cachedPatterns = list(self.patterns)
            self.misses += 1

            if self.match:
                pat, ret = self, self.match.search(data, pos)
            elif self.begin:
                pat, ret = self, self.begin.search(data, pos)
            elif self.include:
                if self.include == "$self" or self.include == "$base":
                    pat, ret = self.syntax._Syntax__rootPattern.cache(data, pos)
                elif self.include.startswith("source"):
                    syntax = Editor().get_syntax_for_scope(self.include)
                    if syntax:
                        pat, ret = syntax._Syntax__rootPattern.cache(data, pos)
                else:
                    key = self.include[1:]
                    if key in self.syntax._Syntax__repo:
                        pat = self.syntax._Syntax__repo[key]
                        pat, ret = pat.cache(data, pos)
            else:
                pat, ret = self.firstMatch(data, pos, self.__cachedPatterns)

            self.__cachedData = data
            self.__cachedMatch = ret
            self.__cachedPat = pat
            return pat, ret

        def dump(self):
            print self.hits, self.misses
            for pat in self.patterns:
                if isinstance(pat, Editor()._Editor__SyntaxPattern):
                    pat.dump()

        def innerApply(self, scope, lastIdx, match, captures):
            scopes = []
            if captures:
                scopes2 = [scope]
                if "0" in captures:
                    scopes2.append(captures["0"].name)

                for i in range(1, len(match.groups())+1):
                    data = match.group(i)
                    if not data:
                        continue

                    if lastIdx < match.start(i):
                        scope = " ".join(scopes2)
                        scopes.append(Editor()._Editor__Scope(scope, sublime.Region(lastIdx, match.start(i))))

                    si = str(i)
                    if si in captures:
                        scopes2.append(captures[si].name)

                    scope = " ".join(scopes2)
                    if lastIdx <= match.start(i):
                        scopes.append(Editor()._Editor__Scope(scope, sublime.Region(match.start(i), match.end(i))))

                    if si in captures:
                        scopes2.pop()

                    lastIdx = match.end(i)
                if lastIdx < match.end():
                    scope = " ".join(scopes2)
                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(lastIdx, match.end())))
            else:
                if lastIdx < match.end():
                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(lastIdx, match.end())))
            return scopes

        def apply(self, data, scope, match):
            scopes = []
            if self.name:
                scope += " %s" % self.name
            if self.match:
                scopes.extend(self.innerApply(scope, match.start(), match, self.captures))
            else:
                scopes.extend(self.innerApply(scope, match.start(), match, self.beginCaptures))

                if self.end:
                    found = False
                    i = scopes[-1].region.end() if len(scopes) else match.start()
                    end = len(data)
                    while i < len(data):
                        endmatch = self.end.search(data, i)
                        if endmatch:
                            end = endmatch.end()
                        else:
                            if not found:
                                # oops.. no end found, set it to the next line
                                end = data.find("\n", i);
                                break
                            else:
                                end = i
                                break

                        if (not endmatch or (endmatch and endmatch.start() != i+1)) and len(self.__cachedPatterns):
                            pattern2, match2 = self.firstMatch(data, i, self.__cachedPatterns)
                            if pattern2 and match2 and \
                                    ((not endmatch and match2.start() < end) or
                                     (endmatch and match2.start() < endmatch.start())):
                                if i != match2.start():
                                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(i, match2.start())))

                                found = True
                                innerScopes = pattern2.apply(data, scope, match2)
                                scopes.extend(innerScopes)
                                i = innerScopes[-1].region.end()
                                continue

                        if endmatch:
                            scopes.extend(self.innerApply(scope, i, endmatch, self.endCaptures))
                            i = end = scopes[-1].region.end()

                        break
            return scopes

        def firstMatch(self, data, pos, patterns):
            # Find the pattern that is the earliest match
            match = None
            startIdx = -1
            pattern = None
            i = 0
            e = Editor()

            while i < len(patterns):
                innermatch = None
                innerPattern, innermatch = patterns[i].cache(data, pos)
                if innermatch:
                    idx = innermatch.start()
                    if startIdx < 0 or startIdx > idx:
                        startIdx = idx
                        match = innermatch
                        pattern = innerPattern
                        # right at the start, we're not going to find a better pattern than this
                        if idx == pos:
                            break
                    i += 1
                else:
                    # No match was found so no point in looking for it again after this point.
                    patterns.pop(i)
            return (pattern, match);

    class __Syntax:

        def __init__(self, name):
            self.file = name
            self.__data = plistlib.readPlist(name)
            self.__repo = {}
            e = Editor()
            if "repository" in self.__data:
                for key in self.__data["repository"]:
                    repo = self.__data["repository"][key]
                    self.__repo[key] = e._Editor__SyntaxPattern(repo, self)

            self.__rootPattern = e._Editor__SyntaxPattern(self.__data, self)
            self.__scopeName = self.__data["scopeName"]

        def name(self):
            return self.__data["name"] if "name" in self.__data else None

        def get_default_scope(self):
            return self.__scopeName

        def extract_scopes(self, data):
            scopes = []
            maxiter = 10000
            i = 0

            while i < len(data) and maxiter > 0:
                maxiter -= 1
                scope = self.__scopeName
                pattern, match = self.__rootPattern.cache(data, i)
                if not match:
                    break
                innerScopes = pattern.apply(data, self.__scopeName, match)
                if match.start() != i:
                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(i, match.start())))
                scopes.extend(innerScopes)
                i = scopes[-1].region.end()
            #self.__rootPattern.dump()
            return scopes


    class __ColorScheme:
        def __init__(self, name):
            self.data = plistlib.readPlist(name)
            self.settings = {}
            for setting in self.data["settings"]:
                if "settings" in setting:
                    name = setting["scope"] if "scope" in setting else "default"
                    self.settings[name] = setting["settings"]
            self.cache = {}

        def getStyleNameForScope(self, scopes):
            key = scopes
            if scopes in self.cache:
                return self.cache[scopes]

            while len(scopes):
                for scope in self.settings:
                    if scopes.endswith(scope):
                        self.cache[key] = scope
                        return scope
                i1 = scopes.rfind(".")
                i2 = scopes.rfind(" ")
                if i1 == i2:
                    break
                scopes = scopes[0:max(i1, i2)]
            # No particular setting found for this style. Apply default
            self.cache[key] = "default"
            return "default"

        def getStyle(self, name):
            return self.settings[name]


    class __Log:
        def __init__(self, editor):
            self.__log = ""
            self.__old = sys.stdout
            self.__editor = editor
            sys.stdout = self
            sys.stderr = self

        def write(self, data):
            self.__log += data
            c = self.__editor.get_console()
            if c:
                 e = c.begin_edit()
                 c.replace(e, sublime.Region(0, c.size()), self.__log)
                 c.end_edit(e)
            if not self.__editor.settings().get("disable_stdout", False):
                self.__old.write(data)

        def flush(self):
            pass

        def dump(self):
            self.__old.write(self.__log)

    def __init__(self):
        self.__log = self.__Log(self)
        start = time.time()
        self.__user_data_dir = appdirs.user_data_dir("lime")
        if not os.path.isdir(self.__user_data_dir):
            os.mkdir(self.__user_data_dir)
        self.__settings = self.__Settings()

        path = sublime.packages_path()
        names = ["Default/Preferences.sublime-settings",
                 "Default/Preferences%s.sublime-settings" % self.platform_settings_name(),
                 "User/Preferences.sublime-settings",
                 "User/Preferences%s.sublime-settings" % self.platform_settings_name()]
        for name in names:
            name = "%s/%s" % (path, name)
            self.__settings._Settings__update(name)

        self.scheme = self.__ColorScheme("%s/../%s" % (sublime.packages_path(), self.__settings.get("color_scheme")))
        self.__windows = []
        syntaxScopes = "%s/syntaxes.cache" % self.__user_data_dir
        if os.path.isfile(syntaxScopes):
            f = gzip.GzipFile(syntaxScopes, "rb")
            self.__syntaxScopes = pickle.load(f)
            self.__syntaxExtensions = pickle.load(f)
            f.close()
        else:
            self.__syntaxScopes, self.__syntaxExtensions = self.__loadsyntaxes()
            f = gzip.GzipFile(syntaxScopes, "wb")
            pickle.dump(self.__syntaxScopes, f)
            pickle.dump(self.__syntaxExtensions, f)
            f.close()
        self.__syntaxCache = {}
        self.__console = None
        print "init took %f ms" % (1000*(time.time()-start))
        self.__tasks = Queue.Queue()
        self.__add_task(self.__load_stuff)

    def get_console(self):
        return self.__console


    def __add_task(self, task, *args):
        self.__tasks.put((task, args))

    def update(self):
        if not self.__tasks.empty():
            try:
                task, args = self.__tasks.get()
                task(*args)
            except:
                traceback.print_exc()
            finally:
                self.__tasks.task_done()
        return not self.__tasks.empty()

    def exit(self, code=0):
        import threading
        for thread in threading.enumerate():
            if threading.current_thread() ==  thread:
                continue

            if thread.isAlive():
                try:
                    # TODO: is there a documented way to do this?
                    thread._Thread__stop()
                except:
                    print(str(thread.getName()) + ' could not be terminated')
        print "Good bye"
        if self.__settings.get("disable_stdout", False):
            self.__log.dump()
        sys.exit(code)

    def get_syntax_file_for_filename(self, name):
        ext = re.search(r"(?<=\.)([^\.]+)$", name)
        ext = ext.group(1) if ext else ""
        if ext in self.__syntaxExtensions:
            return self.__syntaxExtensions[ext]
        return None

    def get_syntax_for_scope(self, name):
        if name in self.__syntaxScopes:
            return self.__get_syntax(self.__syntaxScopes[name])
        return None

    def __get_syntax(self, name):
        name = os.path.abspath(name)
        if name not in self.__syntaxCache:
            print "loading %s" % name
            self.__syntaxCache[name] = self.__Syntax(name)
        return self.__syntaxCache[name]

    def get_syntax(self, name):
        if not name.startswith(sublime.packages_path()):
            name = "%s/../%s" % (sublime.packages_path(), name)
        return self.__get_syntax(name)

    def new_window(self):
        ret = self.__Window()
        if not self.__console:
            self.__console = self.__View(ret)
            self.__console.set_scratch(True)

        self.__windows.append(ret)
        return ret

    def platform_settings_name(self):
        lut = {"osx": " (OSX)", "linux": " (Linux)", "windows": " (Windows)"}
        return lut[sublime.platform()]

    def __load_stuff(self):
        start = time.time()
        keys = []
        path = sublime.packages_path()
        commands = []
        oskeymap = "Default%s.sublime-keymap" % self.platform_settings_name()

        for filename in os.listdir(path):
            filename = "%s/%s" % (path, filename)
            if os.path.isdir(filename):
                for km in ["Default.sublime-keymap", oskeymap]:
                    km = "%s/%s" % (filename, km)
                    if os.path.isfile(km):
                        try:
                            keys.extend(loadjson(km))
                        except:
                            print "Failed to load keymap %s" % km
                            traceback.print_exc()
                cm = "%s/Default.sublime-commands" % filename
                if os.path.isfile(cm):
                    try:
                        commands.extend(loadjson(cm))
                    except:
                        print "Failed to load commands %s" % cm
                        traceback.print_exc()
                for filename2 in os.listdir(filename):
                    if filename2.endswith(".py") and filename2 != "setup.py":
                        filename2 = "%s/%s" % (filename, filename2)
                        self.__add_task(sublime_plugin.reload_plugin, filename2)
        def log_load_time(start):
            print "Loading plugins/commands/keymaps took %f ms" % (1000*(time.time()-start))
        self.__add_task(log_load_time, start)


    def __loadsyntaxes(self):
        path = sublime.packages_path()
        syntaxes = {}
        extensions = {}
        for filename in os.listdir(path):
            filename = "%s/%s" % (path, filename)
            if os.path.isdir(filename):
                for filename2 in os.listdir(filename):
                    filename2 = "%s/%s" % (filename, filename2)
                    if filename2.endswith(".tmLanguage"):
                        try:
                            plist = plistlib.readPlist(filename2)
                            if "scopeName" in plist:
                                data = plist["scopeName"]
                                syntaxes[data] = filename2
                            if "fileTypes" in plist:
                                for f in plist["fileTypes"]:
                                    extensions[f] = filename2
                        except:
                            print "Failed parsing syntax \"%s\"" % filename2
                            traceback.print_exc()
        return syntaxes, extensions

    def settings(self):
        return self.__settings

    def windows(self):
        return self.__windows

    def active_window(self):
        return None

import sublime
import sublime_plugin
