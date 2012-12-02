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

def cssify(item, attrs=[]):
    item = editor.theme().get_class(item, attrs)
    if not item:
        return ""
    ret = ""
    if "color" in item:
        r,g,b = item["color"][0:3]
        ret += "\tcolor: #%s;\n" % ("%02x%02x%02x" % (r,g,b))
    if "fg" in item:
        r,g,b = item["fg"][0:3]
        ret += "\tcolor: #%s;\n" % ("%02x%02x%02x" % (r,g,b))
    if "bg" in item:
        r,g,b = item["bg"][0:3]
        ret += "\tbackground-color: #%s;\n" % ("%02x%02x%02x" % (r,g,b))

    if "layer0.texture" in item:
        offsets = "1"
        if "layer0.inner_margin" in item:
            i = item["layer0.inner_margin"]
            if len(i) == 2:
                i = [i[0], i[1], i[0],i[1]]
            i = [i[1], i[0], i[3], i[2]]
            offsets = " ".join([str(x) for x in i])
        if "content_margin" in item:
            ret += "\tpadding: %spx;\n" % "px ".join([str(x) for x in item["content_margin"]])

        ret += "\tborder-image: url(%s/%s) %s;\n" % (sublime.packages_path(), item["layer0.texture"], offsets)

    return ret

def cssify_theme():

    ret = ""
    ret += "QStatusBar\n{\n"
    ret += cssify("label_control")
    ret += cssify("status_bar")
    ret += "}\n"
    ret += "QScrollArea { border: 0px; }\n"

    ret += "QTabWidget { position: absolute; left: 0; color: yellow; background-color: yellow; }\n"
    ret += "QTabWidget::pane {  color: pink; background-color: pink;  }\n"
    ret += "QTabWidget::tab-bar\n{\n"
    # ret += cssify("tabset_control")
    ret += "subcontrol-position: left;\n"
    ret += "background-color: blue;\n"
    ret += "alignment: left;\n"
    ret += "color: blue;\n"
    ret += "}\n"

    ret += "QTabBar::tab { left: 0px; }"
    ret += "QTabBar::close-button {subcontrol-position: right;}"
    # "subcontrol-origin: border;  position: absolute; left: 0; right: 0; background-color: blue; padding-right: 0; }\n"
    ret += "QTabBar::tab\n{\n"
    ret += cssify("tab_label")
    ret += cssify("tab_control")
    item = editor.theme().get_class("tabset_control")
    ret += "position: absolute; left: 0px;"
    ret += "alignment: left; subcontrol-origin: border; subcontrol-position: left;"
    ret += "border: 0px; margin: 0px; padding: 0px;"
    #ret += "height: %dpx;" % item["tab_height"]
    ret += "}\n"

    ret += "QTabBar::tab:hover\n{\n"
    ret += cssify("tab_control", ["hover"])
    ret += "}\n"

    ret += "QTabBar::tab:selected\n{\n"
    ret += cssify("tab_control", ["selected"])
    ret += "}\n"
    return str(ret)

def transform_scopes(view):
    scopes = [(editor.scheme().getStyleNameForScope(scope.name).replace(".", "_dot_"), scope.region) for scope in view._View__scopes]
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
        output = output.replace(" ", "&nbsp;").replace("\n", "<br>\n")
        od += "<span class=\"%s\">%s</span>" % (name, output)
    od += "</body>"
    return od


def create_stylesheet():
    css = ""
    lut = {"foreground": "color", "background": "background"}
    for name in editor.scheme().settings:
        settings = editor.scheme().settings[name]
        config = ""
        for key in settings:
            if key in lut:
                if len(config):
                    config += " "
                config += "%s: %s;" % (lut[key], settings[key])
        name = name.replace(".", "_dot_")

        css += ".%s { %s }\n" % (name, config)

    return "<style>\n%s</style>" % css

def background_color():
    col = editor.scheme().settings["default"]["background"]
    match = re.match(r"#([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})([A-Fa-f0-9]{2})", col)
    return [int(a, 16) for a in match.groups()]


print background_color()
