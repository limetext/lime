try:
    import traceback
    import sublime
    print("new file")
    v = sublime.test_window.new_file()
    print("running command")
    v.run_command("test_text")
    print("command ran")
    assert v.substr(sublime.Region(0, v.size())) == "hello"
    v.run_command("undo")
    print(v.sel()[0])
    assert v.sel()[0] == (0, 0)
    v = sublime.test_window.active_view()
    sublime.test_window.run_command("test_window")
    assert v.substr(sublime.Region(0, v.size())) == "window hello"
    assert sublime.CLASS_PUNCTUATION_START == 4
    assert sublime.CLASS_OPENING_PARENTHESIS == 4096
except:
    traceback.print_exc()
    raise
