import traceback
try:
    import unittest
    a = unittest.defaultTestLoader
except:
    traceback.print_exc()
    raise
