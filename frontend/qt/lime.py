from curses import *
import curses
import sys
import os.path
sys.path.append("%s/../../backend" % os.path.dirname(os.path.abspath(__file__)))
import backend
import re
import traceback
import sublime
import time
import sublime_plugin

editor = backend.Editor()

def transform_scopes(view):
    scopes = [(editor.scheme.getStyleNameForScope(scope.name).replace(".", "_dot_"), scope.region) for scope in view._View__scopes]
    # Merge scopes that have the same style
    colorscopes = [scopes[0]]
    for n2,r2 in scopes[1:]:
        n1,r1 = colorscopes.pop()

        if n1 == n2:
            r1 = r1.cover(r2)
            colorscopes.append((n1, r1))
        else:
            colorscopes.append((n1, r1))
            colorscopes.append((n2, r2))

    data = view.substr(sublime.Region(0, view.size()))
    od = "<body class=\"default\">"
    for name, region in colorscopes:
        output = data[region.begin():region.end()]
        output = output.replace(" ", "&nbsp;").replace("\n", "<br>")
        od += "<span class=\"%s\">%s</span>" % (name, output)
    od += "</body>"
    return od


def create_stylesheet():
    css = ""
    lut = {"foreground": "color", "background": "background"}
    for name in editor.scheme.settings:
        settings = editor.scheme.settings[name]
        config = ""
        for key in settings:
            if key in lut:
                if len(config):
                    config += " "
                config += "%s: %s;" % (lut[key], settings[key])
        name = name.replace(".", "_dot_")

        css += ".%s { %s }\n" % (name, config)

    return "<style>\n%s</style>" % css

