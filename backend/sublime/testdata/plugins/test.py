import sublime, sublime_plugin

class TestText(sublime_plugin.TextCommand):
    def run(self, edit):
        print("my view's id is: %d" % self.view.id())
        self.view.insert(edit, 0, "hello")

class TestWindow(sublime_plugin.WindowCommand):
    def run(self):
        print("my window's id is %d" % self.window.id())
        v = self.window.active_view()
        e = v.begin_edit()
        v.insert(e, 0, "window hello")
        v.end_edit(e)
