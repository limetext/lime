# coding=utf-8
import sys
import traceback
try:
    import sublime
    v = sublime.test_window.new_file()
    assert v.id() != sublime.test_window.id()
    assert sublime.test_window.id() == v.window().id()
    assert v.size() == 0
    e = v.begin_edit()
    v.insert(e, 0, "hellå world")
    v.end_edit(e)
    assert v.substr(sublime.Region(0, v.size())) == "hellå world"
    e = v.begin_edit()
    v.insert(e, 0, """abrakadabra
simsalabim
hocus pocus
""")
    v.end_edit(e)
    assert v.rowcol(20) == (1, 8)
    assert v.rowcol(25) == (2, 2)

    assert len(v.sel()) == 1
    assert len(list(v.sel())) == 1
    assert v.settings().get("test", "hello") == "hello"
    v.settings().set("test", 10)
    assert v.settings().get("test") == 10
    assert v.sel()[0] == (46, 46)
    v.run_command("move", {"by": "characters", "forward": False})
    assert v.sel()[0] == (45, 45)
    v.run_command("move", {"by": "characters", "forward": True})
    assert v.sel()[0] == (46, 46)

    v2 = sublime.test_window.new_file()
    e = v2.begin_edit()
    v2.insert(e, 0, """one word { another word }
line

after empty line""")
    v2.end_edit(e)
    # Expected results validated in Sublime
    assert v2.find_by_class(1, True, sublime.CLASS_WORD_START) == sublime.Region(4, 4)
    assert v2.find_by_class(1, True, sublime.CLASS_PUNCTUATION_START) == sublime.Region(9, 9)
    assert v2.expand_by_class(sublime.Region(5, 6),
    	sublime.CLASS_WORD_START | sublime.CLASS_WORD_END) == sublime.Region(4, 8)
    assert v2.expand_by_class(sublime.Region(11, 12),
    	sublime.CLASS_PUNCTUATION_START | sublime.CLASS_PUNCTUATION_END) == sublime.Region(10, 24)
    assert v2.expand_by_class(sublime.Region(5, 6), sublime.CLASS_EMPTY_LINE) == sublime.Region(0, 31)

except:
    print(sys.exc_info()[1])
    traceback.print_exc()
    raise
