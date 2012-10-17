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

class Scope:
    def __init__(self, name, region):
        self.name   = name
        self.region = region

    def __repr__(self):
        return "%s %s" % (self.name, self.region)

class SyntaxPattern:
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
                match = SyntaxPattern.HackRegex.HackMatchObject(match, idx)
            return match

    def compile(self, data, var):
        try:
            data = data[var] if var in data else None
            if data:
                if "\\G" in data:
                    data = SyntaxPattern.HackRegex(data)
                else:
                    data = re.compile(data)
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

        if "patterns" in data:
            for pattern in data["patterns"]:
                if "include" in pattern:
                    pattern = pattern["include"][1:]
                    pattern = syntax._Syntax__repo[pattern] if pattern in syntax._Syntax__repo else None
                    if pattern:
                        for pat in pattern:
                            self.patterns.append(pat)
                else:
                    self.patterns.append(SyntaxPattern(pattern, syntax))
        self.name = data["name"] if "name" in data else ""

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
                if lastIdx != match.start(i):
                    scope = " ".join(scopes2)
                    scopes.append(Scope(scope, sublime.Region(lastIdx, match.start(i))))

                si = str(i)
                if si in captures:
                    scopes2.append(captures[si].name)

                scope = " ".join(scopes2)
                scopes.append(Scope(scope, sublime.Region(match.start(i), match.end(i))))

                if si in captures:
                    scopes2.pop()

                lastIdx = match.end(i)
            if lastIdx != match.end():
                scope = " ".join(scopes2)
                scopes.append(Scope(scope, sublime.Region(lastIdx, match.end())))
        else:
            if lastIdx != match.end():
                scopes.append(Scope(scope, sublime.Region(lastIdx, match.end())))
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
                i = match.end()
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

                    if self.patterns:
                        pattern2, match2 = self.firstMatch(data, i, self.patterns)
                        if pattern2 and match2 and \
                                ((not endmatch and match2.start() < end) or
                                 (endmatch and match2.start() < endmatch.start())):
                            if i != match2.start():
                                scopes.append(Scope(scope, sublime.Region(i, match2.start())))

                            found = True
                            innerScopes = pattern2.apply(data, scope, match2)
                            scopes.extend(innerScopes)
                            i = innerScopes[-1].region.end()
                            continue

                    if endmatch:
                        scopes.extend(self.innerApply(scope, i, endmatch, self.endCaptures))
                        i = end = endmatch.end()

                    break
        return scopes

    def firstMatch(self, data, pos, patterns, cache=None, remove=False):
        # Find the pattern that is the earliest match
        match = None
        startIdx = -1
        pattern = None
        i = 0
        while i < len(patterns):
            innerPattern = patterns[i]
            innermatch = None
            if cache and cache[i]:
                innermatch = cache[i]
            else:
                if innerPattern.match:
                    innermatch = innerPattern.match.search(data, pos)
                elif innerPattern.begin:
                    innermatch = innerPattern.begin.search(data, pos)
            if cache:
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

class Syntax:

    def __init__(self, name):
        self.__data = plistlib.readPlist(name)
        self.__repo = {}
        if "repository" in self.__data:
            for key in self.__data["repository"]:
                repo = self.__data["repository"][key]["patterns"] or []
                if len(repo):
                    self.__repo[key] = []
                for i in range(len(repo)):
                    self.__repo[key].append(SyntaxPattern(repo[i], self))

        self.__rootPattern = SyntaxPattern(self.__data, self)
        self.__scopeName = self.__data["scopeName"]

    def __firstMatch(self, data, pos, patterns, cache, remove):
        return self.__rootPattern.firstMatch(data, pos, patterns, cache, remove)

    def extract_scopes(self, data):
        scopes = []
        maxiter = 10000
        i = 0
        cache = [None for a in self.__rootPattern.patterns]
        patterns = list(self.__rootPattern.patterns)
        while i < len(data) and len(patterns) and maxiter > 0:
            maxiter -= 1
            scope = self.__scopeName
            pattern, match = self.__firstMatch(data, i, patterns, cache, True)
            if not match:
                break
            innerScopes = pattern.apply(data, self.__scopeName, match)
            if match.start() != i:
                scopes.append(Scope(scope, sublime.Region(i, match.start())))
            scopes.extend(innerScopes)

            i = scopes[-1].region.end()

            for j in range(len(cache)):
                if cache[j] and (cache[j].end() < i or cache[j].start() < i):
                    cache[j] = None
        return scopes


class ColorScheme:
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

def singleton(cls):
    instances = {}
    def getinstance():
        if cls not in instances:
            instances[cls] = cls()
        return instances[cls]
    return getinstance



@singleton
class Editor:
    class Window:
        def __init__(self):
            self.__views = {}
            e = Editor()
            self.__settings = e.Settings(e._Editor__settings)
            # TODO: project settings

        def open_file(self, name, flags=0):
            # TOOD: handle sublime.ENCODED_POSITION
            if name not in self.__views:
                e = Editor()
                self.__views[name] = e.View(self, name)
            return self.__views[name]



    class View:
        def __init__(self, window, name=None):
            if name and os.path.isfile(name):
                f = open(name)
                self.__buffer = f.read()
                f.close()
            self.__window = window
            try:
                e = Editor()
                self.__settings = e.Settings(window._Window__settings)
                # TODO: dynamically detect syntax
                # TODO: apply syntax specific settings
            except:
                traceback.print_exc()

        def window(self):
            return self.__window

        def settings(self):
            return self.__settings

        def size(self):
            return len(self.__buffer)

        def substr(self, region):
            return self.__buffer[region.begin():region.end()]

    class Settings:
        def __init__(self, other=None):
            if other:
                self.__values = dict(other.__values)
            else:
                self.__values = {}

        def get(self, name, default=None):
            if name in self.__values:
                return self.__values[name]
            return default

        def set(self, name, value):
            self.__values[name] = value

        def erase(self, name):
            if name in self.__values:
                del self.__values[name]

        def __update(self, other):
            if isinstance(other, str):
                if os.path.isfile(other):
                    self.__values.update(loadjson(other))
            else:
                self.__values.update(other.__values)

    def __init__(self):
        self.__loadkeymaps()
        self.__settings = self.Settings()

        path = sublime.packages_path()
        names = ["Default/Preferences.sublime-settings",
                 "Default/Preferences%s.sublime-settings" % self.platform_settings_name(),
                 "User/Preferences.sublime-settings",
                 "User/Preferences%s.sublime-settings" % self.platform_settings_name()]
        for name in names:
            name = "%s/%s" % (path, name)
            self.__settings._Settings__update(name)


        self.scheme = ColorScheme("%s/../%s" % (sublime.packages_path(), self.__settings.get("color_scheme")))
        self.__windows = []

    def new_window(self):
        ret = self.Window()
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
