import sys
import traceback
try:
	import sublime

	r = sublime.Region()
	assert r.a == 0 and r.b == 0
	r = sublime.Region(1)
	assert r.a == 1 and r.b == 0
	r = sublime.Region(1, 2)
	assert r.a == 1 and r.b == 2

	ok = False
	try:
		r = sublime.Region(1, 2, 3)
	except:
		ok = True
	assert ok

	r = sublime.Region(1, 3)
	assert r.contains(1) and not r.contains(4)
	r2 = sublime.Region(5, 10)

	r3 = r.cover(r2)
	assert r3.a == 1 and r3.b == 10

	ok = False
	try:
		r4 = r.cover(4)
	except:
		ok = True
	assert ok

	r3 = sublime.Region(3, 2)
	assert r3.begin() == 2 and r3.end() == 3
except:
	print(sys.exc_info()[1])
	traceback.print_exc()
	raise
