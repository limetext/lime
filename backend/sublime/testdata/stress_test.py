import sys
import traceback
import types
try:
    import sublime
    v = sublime.test_window.new_file()
    for i in range(10000):
        e = v.begin_edit()
        v.insert(e, 0, "hello world")
        v.erase(e, sublime.Region(0, 11))
        v.insert(e, 0, "hello world")
        v.end_edit(e)
        assert v.substr(sublime.Region(0, v.size())) == "hello world"
        v.run_command("undo")
        assert v.sel()[0] == (0, 0)

except:
    print sys.exc_info()[1]
    traceback.print_exc()
    raise
