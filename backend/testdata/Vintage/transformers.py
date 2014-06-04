import sublime
import sublime_plugin

from Vintageous.vi.constants import regions_transformer
from Vintageous.vi.utils import back_one_char
from Vintageous.vi.utils import forward_one_char
from Vintageous.vi.utils import is_at_bol
from Vintageous.vi.utils import is_at_eol
from Vintageous.vi.utils import is_at_hard_eol
from Vintageous.vi.utils import is_line_empty
from Vintageous.vi.utils import is_on_empty_line
from Vintageous.vi.utils import next_non_white_space_char


class ClipEndToLine(sublime_plugin.TextCommand):
    def run(self, edit):
        def f(view, s):
            if not is_on_empty_line(self.view, s) and is_at_eol(self.view, s):
                return back_one_char(s)
            else:
                return s

        regions_transformer(self.view, f)


class DontStayOnEolForward(sublime_plugin.TextCommand):
    def run(self, edit, **kwargs):
        def f(view, s):
            if is_at_eol(self.view, s):
                return forward_one_char(s)
            else:
                return s

        regions_transformer(self.view, f)


class DontStayOnEolBackward(sublime_plugin.TextCommand):
    def run(self, edit, **kwargs):
        def f(view, s):
            if is_at_eol(self.view, s) and not self.view.line(s.b).empty():
                return back_one_char(s)
            else:
                return s

        regions_transformer(self.view, f)


class _vi_d_post_action(sublime_plugin.TextCommand):
    def run(self, edit, **kwargs):
        def f(view, s):
            if is_at_eol(self.view, s) and not self.view.line(s.b).empty():
                s = back_one_char(s)
            # s = next_non_white_space_char(self.view, s.b)
            return s

        regions_transformer(self.view, f)


class DontOvershootLineLeft(sublime_plugin.TextCommand):
    def run(self, edit, **kwargs):
        def f(view, s):
            if view.size() > 0 and is_at_eol(self.view, s):
                return forward_one_char(s)
            else:
                return s

        regions_transformer(self.view, f)