import re
import plistlib
import traceback
import sys
import json
import os
import os.path

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

    class __View:
        def __init__(self, window, name=None):
            if name and os.path.isfile(name):
                f = open(name)
                self.__buffer = f.read()
                f.close()
            self.__window = window
            try:
                e = Editor()
                self.__settings = e._Editor__Settings(window._Window__settings)
                self.__syntax = None
                self.__settings.add_on_change("__lime_View", self.__settings_changed)
                # TODO: dynamically detect syntax
                self.__settings.set("syntax", "Packages/JavaScript/JavaScript.tmLanguage")
            except:
                traceback.print_exc()

        def __settings_changed(self):
            syntax = self.__settings.get("syntax")
            if self.__syntax == None or self.__syntax.file != syntax:
                e = Editor()
                self.__syntax = e.get_syntax(syntax)
                print self.__syntax.file
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

    class __Settings:
        def __init__(self, other=None):
            if other:
                self.__values = dict(other.__values)
            else:
                self.__values = {}
            self.__on_change = {}

        def add_on_change(self, key, callback):
            if key not in self.__on_change:
                self.__on_change[key] = []
            self.__on_change[key].append(callback)

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
                    for cb in self.__on_change[key]:
                        cb()
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
                    regex = re.compile(r"(\(\?[iLmsux]+):")
                    match = regex.search(data)
                    while match:
                        data = "%s(?:%s)%s" % (data[:match.start()], match.group(1), data[match.end():])
                        match = regex.search(data)

                    regex = re.compile(r"(\(\?<)")
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

            if "patterns" in data:
                e = Editor()
                for pattern in data["patterns"]:
                    if "include" in pattern:
                        pattern = pattern["include"]
                        self.patterns.append(pattern)
                    else:
                        self.patterns.append(e._Editor__SyntaxPattern(pattern, syntax))
            self.name = data["name"] if "name" in data else None
            if self.name == None and "contentName" in data:
                self.name = data["contentName"]


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
                    verify(scopes)

                    si = str(i)
                    if si in captures:
                        scopes2.append(captures[si].name)

                    scope = " ".join(scopes2)
                    if lastIdx <= match.start(i):
                        scopes.append(Editor()._Editor__Scope(scope, sublime.Region(match.start(i), match.end(i))))
                    verify(scopes)

                    if si in captures:
                        scopes2.pop()

                    lastIdx = match.end(i)
                verify(scopes)
                if lastIdx < match.end():
                    scope = " ".join(scopes2)
                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(lastIdx, match.end())))
                    verify(scopes)
            else:
                if lastIdx < match.end():
                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(lastIdx, match.end())))
                    verify(scopes)
            verify(scopes)
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
                    patterns = list(self.patterns)
                    cache = [None for i in patterns]
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
                            else:
                                end = i
                                break

                        if len(patterns):
                            pattern2, match2 = self.firstMatch(data, i, patterns, cache, True)
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

        def firstMatch(self, data, pos, patterns, cache=None, remove=False):
            # Find the pattern that is the earliest match
            match = None
            startIdx = -1
            pattern = None
            i = 0
            e = Editor()

            if cache:
                for j in range(len(cache)):
                    if cache[j] and (cache[j].end() < pos or cache[j].start() < pos):
                        cache[j] = None

            while i < len(patterns):
                syntaxMatch = False
                innerPattern = patterns[i]
                innermatch = None
                if cache and cache[i]:
                    innermatch = cache[i]
                else:
                    syntaxMatch = not isinstance(innerPattern, e._Editor__SyntaxPattern)
                    if syntaxMatch:
                        if innerPattern == "$self":
                            innerPattern, innermatch = self.syntax.recurse(data, pos)
                        else:
                            key = innerPattern[1:]
                            if key in self.syntax._Syntax__repo:
                                pat = self.syntax._Syntax__repo[key]
                                pats = []
                                if pat.match or pat.begin:
                                    pats = [pat]
                                else:
                                    pats = pat.patterns
                                innerPattern, innermatch = pat.firstMatch(data, pos, pats)
                            else:
                                innerPattern, innermatch = None, None
                    else:
                        if innerPattern.match:
                            innermatch = innerPattern.match.search(data, pos)
                        elif innerPattern.begin:
                            innermatch = innerPattern.begin.search(data, pos)
                if cache and not syntaxMatch:
                    cache[i] = innermatch
                if innermatch:
                    idx = innermatch.start()
                    if startIdx < 0 or startIdx > idx:
                        startIdx = idx
                        match = innermatch
                        pattern = innerPattern
                if remove and innermatch == None:
                    # No match was found and we've indicated that the pattern can be removed
                    # if that is the case (ie if it wasn't found, it's never going to be found,
                    # so no point in looking for it again after this point).
                    patterns.pop(i)
                    cache.pop(i)
                else:
                    i += 1
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
            self.__recurse = False

        def name(self):
            return self.__data["name"] if "name" in self.__data else None

        def get_default_scope(self):
            return self.__scopeName

        def firstMatch(self, data, pos, patterns, cache, remove):
            return self.__rootPattern.firstMatch(data, pos, patterns, cache, remove)


        def recurse(self, data, pos):
            if self.__recurse:
                return None, None
            self.__recurse = True
            try:
                return self.firstMatch(data, pos, self.__rootPattern.patterns, None, False)
            finally:
                self.__recurse = False

        def extract_scopes(self, data):
            scopes = []
            maxiter = 10000
            i = 0
            cache = [None for a in self.__rootPattern.patterns]
            patterns = list(self.__rootPattern.patterns)

            while i < len(data) and len(patterns) and maxiter > 0:
                maxiter -= 1
                scope = self.__scopeName
                pattern, match = self.firstMatch(data, i, patterns, cache, True)
                if not match:
                    break
                innerScopes = pattern.apply(data, self.__scopeName, match)
                if match.start() != i:
                    scopes.append(Editor()._Editor__Scope(scope, sublime.Region(i, match.start())))
                verify(innerScopes)
                scopes.extend(innerScopes)
                i = scopes[-1].region.end()

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

    def __init__(self):
        self.__loadkeymaps()
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

    def get_syntax(self, name):
        return self.__Syntax("%s/../%s" % (sublime.packages_path(), name))

    def new_window(self):
        ret = self.__Window()
        self.__windows.append(ret)
        return ret

    def platform_settings_name(self):
        lut = {"osx": " (OSX)", "linux": " (Linux)", "windows": " (Windows)"}
        return lut[sublime.platform()]

    def __loadkeymaps(self):
        keys = []
        path = sublime.packages_path()
        oskeymap = "Default%s.sublime-keymap" % self.platform_settings_name()
        for filename in os.listdir(path):
            filename = "%s/%s" % (path, filename)
            if os.path.isdir(filename):
                for km in ["Default.sublime-keymap", oskeymap]:
                    km = "%s/%s" % (filename, km)
                    if os.path.isfile(km):
                            keys.extend(loadjson(km))

    def windows(self):
        return self.__windows

    def active_window(self):
        return None

import sublime
