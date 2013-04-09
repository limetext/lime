try:
    import traceback
    import sublime
    v = sublime.test_window.new_file()
    v.run_command("test_text")
    assert v.substr(sublime.Region(0, v.size())) == "hello"
    v.run_command("undo")
    assert v.sel()[0] == (0, 0)
    v = sublime.test_window.active_view()
    sublime.test_window.run_command("test_window")
    assert v.substr(sublime.Region(0, v.size())) == "window hello"
except:
    traceback.print_exc()
    raise
