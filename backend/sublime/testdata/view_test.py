import sys
import traceback
try:
	import sublime

	v = sublime.test_window.new_view()
	print v, sublime.test_window, v.window()
	print v.id(), sublime.test_window.id(), v.window().id()
except:
	print sys.exc_info()[1]
	traceback.print_exc()
	raise
