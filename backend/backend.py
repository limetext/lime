import sublime
import re
import plistlib
import traceback
import sys

class Scope:
    def __init__(self, name, region):
        self.name   = name
        self.region = region

    def __repr__(self):
        return "%s %s" % (self.name, self.region)

class SyntaxPattern:
    def compile(self, data, var):
        try:
            setattr(self, var, re.compile(data[var]) if var in data else None)
        except:
            setattr(self, var, None)
            print "Failed to compile regex: \"%s\" - %s" % (data[var], sys.exc_info()[1])

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
                    pattern = syntax.repo[pattern] if pattern in syntax.repo else []

                    for pat in pattern:
                        self.patterns.append(pat)
                else:
                    self.patterns.append(SyntaxPattern(pattern, syntax))
        self.name = data["name"] if "name" in data else None

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
                            scopes.extend(pattern2.apply(data, scope, match2))
                            i = match2.end()
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
                if (startIdx < 0 or startIdx > idx) and idx != innermatch.end():
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
        self.data = plistlib.readPlist(name)
        self.repo = {}
        for key in self.data["repository"]:
            repo = self.data["repository"][key]["patterns"] or []
            if len(repo):
                self.repo[key] = []
            for i in range(len(repo)):
                self.repo[key].append(SyntaxPattern(repo[i], self))

        self.rootPattern = SyntaxPattern(self.data, self)
        self.scopeName = self.data["scopeName"]

    def firstMatch(self, data, pos, patterns, cache, remove):
        return self.rootPattern.firstMatch(data, pos, patterns, cache, remove)

    def extract_scopes(self, data):
        scopes = []
        maxiter = 10000
        i = 0
        cache = [None for a in self.rootPattern.patterns]
        patterns = list(self.rootPattern.patterns)
        while i < len(data) and len(patterns) and maxiter > 0:
            maxiter -= 1
            scope = self.scopeName
            pattern, match = self.firstMatch(data, i, patterns, cache, True)
            if not match:
                break
            innerScopes = pattern.apply(data, self.scopeName, match)
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




