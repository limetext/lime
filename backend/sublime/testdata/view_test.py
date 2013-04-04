import sys
import traceback
try:
	import sublime

	v = sublime.test_window.new_view()
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
	assert v.row_col(20) == (2, 9)
	assert v.row_col(25) == (3, 3)


	assert v.settings().get("test", "hello") == "hello"
	v.settings().set("test", 10)
	assert v.settings().get("test") == 10
except:
	print sys.exc_info()[1]
	traceback.print_exc()
	raise
