import sys
import traceback
try:
	import sublime

	r = sublime.Region()
	r.a, r.b = 1, 3
	assert r.contains(1) and not r.contains(4)
	r2 = sublime.Region()
	r2.a, r2.b = 5, 10

	r3 = r.cover(r2)
	assert r3.a == 1 and r3.b == 10

	ok = False
	try:
		r4 = r.cover(4)
	except:
		ok = True
	assert ok
except:
	print sys.exc_info()[1]
	traceback.print_exc()
