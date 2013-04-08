import sys
import traceback
import types
try:
	import sublime
	v = sublime.test_window.new_file()
	assert v.id() != sublime.test_window.id()
	assert sublime.test_window.id() == v.window().id()
	assert v.size() == 0
	e = v.begin_edit()
	v.insert(e, 0, "hello world")
	v.end_edit(e)
	assert v.substr(sublime.Region(0, v.size())) == "hello world"
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
	v.run_command("move", {"by":"characters", "forward": True})
	assert v.sel()[0] == (47, 47)
	v.run_command("move", {"by":"characters", "forward": False})
	assert v.sel()[0] == (46, 46)
except:
	print sys.exc_info()[1]
	traceback.print_exc()
	raise
