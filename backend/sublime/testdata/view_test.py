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
except:
	print sys.exc_info()[1]
	traceback.print_exc()
	raise
