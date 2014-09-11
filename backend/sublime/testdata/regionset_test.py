import sys
import traceback
try:
    import sublime

    ok = False
    try:
        r = sublime.RegionSet()
    except:
        ok = True
    assert ok
except:
    print(sys.exc_info()[1])
    traceback.print_exc()
    raise
