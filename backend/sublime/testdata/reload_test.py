try:
    import traceback
    import sublime
    print("Testing plugin reload")
    print("new file")
    v = sublime.test_window.new_file()
    print("running command")
    v.run_command("test_toxt")
    print("command ran")
    assert v.substr(sublime.Region(0, v.size())) == "Tada"
except:
    traceback.print_exc()
    raise
