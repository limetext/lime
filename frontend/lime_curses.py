from curses import *
import curses
import sys
import os.path
sys.path.append("%s/../backend" % os.path.dirname(os.path.abspath(__file__)))
import backend
import re
import traceback
import sublime

editor = backend.Editor()
wnd = editor.new_window()
view = wnd.open_file("../lime.js")
data = view.substr(sublime.Region(0, view.size()))

syntax = backend.Syntax("%s/../%s" % (sublime.packages_path(), view.settings().get("syntax", "Packages/JavaScript/JavaScript.tmLanguage")))

scopes = syntax.extract_scopes(data)
log = ""

color_start = 30
colors = {}
def add_color(hex):
    global color_start
    if hex in colors:
        return colors[hex]
    match = re.match("\#?([a-fA-F0-9]{2})([a-fA-F0-9]{2})([a-fA-F0-9]{2})", hex)
    ret = color_start

    if match:
        r, g, b = [int(1000*(int(a, 16)/255.0)) for a in match.groups()]
        init_color(color_start, r, g, b)
        colors[hex] = color_start
        color_start += 1
    return ret

color_pairs = 4
pair_lut = {}

def get_color(name):
    global color_pairs
    if name in pair_lut:
        return pair_lut[name]

    settings = editor.scheme.settings[name]
    fg = add_color(settings["foreground"]) if "foreground" in settings else 0
    bg = add_color(settings["background"]) if "background" in settings else 0
    if bg == 0 and "background" in editor.scheme.settings["default"]:
        bg = add_color(editor.scheme.settings["default"]["background"])
    init_pair(color_pairs, fg, bg)
    pair_lut[name] = color_pairs
    ret = color_pairs
    color_pairs += 1
    return ret

def clear():
    lines, columns = stdscr.getmaxyx()
    for i in range(lines):
        stdscr.addstr(i, 0, " " * (columns-1), color_pair(get_color("default")))
    stdscr.move(0, 0)

try:
    stdscr = initscr()


    raw()
    noecho()
    if can_change_color():
        start_color()

    refresh = True

    inv_fg = add_color(editor.scheme.settings["default"]["invisibles"])
    inv_bg = add_color(editor.scheme.settings["default"]["background"])
    init_pair(color_pairs, inv_fg, inv_bg)
    inv_pair = color_pairs
    color_pairs += 1
    inv_regex = re.compile(r"(\n|[\t ]+)")

    while True:
        try:
            if refresh:
                clear()
                stdscr.move(0, 0)

                line = 0
                for scope in scopes:
                    output = data[scope.region.begin():scope.region.end()]
                    color = get_color(editor.scheme.getStyleNameForScope(scope.name))

                    if "\n" in output:
                        output = output.split("\n")
                        if stdscr.getyx()[0] + len(output) >= stdscr.getmaxyx()[0]:
                            break
                        for l in output[:-1]:
                            stdscr.addstr(l, color_pair(color))
                            line += 1
                            stdscr.move(line, 0)
                        output = output[-1]
                    stdscr.addstr(output, color_pair(color))
                line = 0
                for match in inv_regex.finditer(data):
                    l = match.group(1)
                    add = l.count("\n")
                    if line + add >= stdscr.getmaxyx()[0]:
                        break
                    xoff = len(data[:match.start()].split("\n")[-1])

                    l = l.replace(" ", ".").replace("\n", "$").replace("\t", ">---")
                    try:
                        stdscr.addstr(line, xoff, l, color_pair(inv_pair))
                    except:
                        break
                    line += add
                    if line >= stdscr.getmaxyx()[0]:
                        break


                stdscr.refresh()
                refresh = False

            rawch = stdscr.getch()
            if rawch == KEY_RESIZE:
                refresh = True
            ch = keyname(rawch)
            log += "%s\n" % ch
            if ch == "^C":
                break
        except:
            traceback.print_exc()

finally:
    endwin()
    print log
