import sublime
import sublime_plugin

from Vintageous.state import IrreversibleTextCommand
from Vintageous.state import VintageState
from Vintageous.vi import utils
from Vintageous.vi.constants import _MODE_INTERNAL_NORMAL
from Vintageous.vi.constants import MODE_INSERT
from Vintageous.vi.constants import MODE_NORMAL
from Vintageous.vi.constants import MODE_VISUAL
from Vintageous.vi.constants import MODE_VISUAL_LINE
from Vintageous.vi.constants import regions_transformer
from Vintageous.vi.registers import REG_EXPRESSION


class ViEditAtEol(sublime_plugin.TextCommand):
    def run(self, edit, extend=False):
        state = VintageState(self.view)
        state.enter_insert_mode()

        self.view.run_command('collapse_to_direction')

        sels = list(self.view.sel())
        self.view.sel().clear()

        new_sels = []
        for s in sels:
            hard_eol = self.view.line(s.b).end()
            new_sels.append(sublime.Region(hard_eol, hard_eol))

        for s in new_sels:
            self.view.sel().add(s)


class ViEditAfterCaret(sublime_plugin.TextCommand):
    def run(self, edit, extend=False):
        state = VintageState(self.view)
        state.enter_insert_mode()

        visual = self.view.has_non_empty_selection_region()

        sels = list(self.view.sel())
        self.view.sel().clear()

        new_sels = []
        for s in sels:
            if visual:
                new_sels.append(sublime.Region(s.end(), s.end()))
            else:
                if not utils.is_at_eol(self.view, s):
                    new_sels.append(sublime.Region(s.end() + 1, s.end() + 1))
                else:
                    new_sels.append(sublime.Region(s.end(), s.end()))

        for s in new_sels:
            self.view.sel().add(s)


class _vi_big_i(sublime_plugin.TextCommand):
    def run(self, edit, extend=False):
        def f(view, s):
            line = view.line(s.b)
            pt = utils.next_non_white_space_char(view, line.a)
            return sublime.Region(pt, pt)

        state = VintageState(self.view)
        state.enter_insert_mode()

        regions_transformer(self.view, f)


class ViPaste(sublime_plugin.TextCommand):
    def run(self, edit, register=None, count=1):
        state = VintageState(self.view)

        if register:
            fragments = state.registers[register]
        else:
            # TODO: There should be a simpler way of getting the unnamed register's content.
            fragments = state.registers['"']
            if not fragments:
                print("Vintageous: Nothing in register \".")
                # XXX: This won't ever be printed because it will be overwritten by other status
                # messages printed right after this one.
                sublime.status_message("Vintageous: Nothing in register \".")
                return

        sels = list(self.view.sel())

        if len(sels) == len(fragments):
            sel_frag = zip(sels, fragments)
        else:
            sel_frag = zip(sels, [fragments[0],] * len(sels))

        offset = 0
        for s, text in sel_frag:
            text = self.prepare_fragment(text)
            if text.startswith('\n'):
                if utils.is_at_eol(self.view, s) or utils.is_at_bol(self.view, s):
                    self.paste_all(edit, s, self.view.line(s.b).b, text, count)
                else:
                    self.paste_all(edit, s, self.view.line(s.b - 1).b, text, count)
            else:
                # XXX: Refactor this whole class. It's getting out of hand.
                if self.view.substr(s.b) == '\n':
                    self.paste_all(edit, s, s.b + offset, text, count)
                else:
                    self.paste_all(edit, s, s.b + offset + 1, text, count)
                offset += len(text) * count

    def prepare_fragment(self, text):
        if text.endswith('\n') and text != '\n':
            text = '\n' + text[0:-1]
        return text

    # TODO: Improve this signature.
    def paste_all(self, edit, sel, at, text, count):
        state = VintageState(self.view)
        if state.mode not in (MODE_VISUAL, MODE_VISUAL_LINE):
            # TODO: generate string first, then insert?
            # Make sure we can paste at EOF.
            at = at if at <= self.view.size() else self.view.size()
            for x in range(count):
                self.view.insert(edit, at, text)
        else:
            if text.startswith('\n'):
                text = text * count
                if not text.endswith('\n'):
                    text = text + '\n'
            else:
                text = text * count

            if state.mode == MODE_VISUAL_LINE:
                if text.startswith('\n'):
                    text = text[1:]

            self.view.replace(edit, sel, text)


class ViPasteBefore(sublime_plugin.TextCommand):
    def run(self, edit, register=None, count=1):
        state = VintageState(self.view)

        if register:
            fragments = state.registers[register]
        else:
            # TODO: There should be a simpler way of getting the unnamed register's content.
            fragments = state.registers['"']

        sels = list(self.view.sel())

        if len(sels) == len(fragments):
            sel_frag = zip(sels, fragments)
        else:
            sel_frag = zip(sels, [fragments[0],] * len(sels))

        offset = 0
        for s, text in sel_frag:
            if text.endswith('\n'):
                if utils.is_at_eol(self.view, s) or utils.is_at_bol(self.view, s):
                    self.paste_all(edit, s, self.view.line(s.b).a, text, count)
                else:
                    self.paste_all(edit, s, self.view.line(s.b - 1).a, text, count)
            else:
                self.paste_all(edit, s, s.b + offset, text, count)
                offset += len(text) * count

    def paste_all(self, edit, sel, at, text, count):
        # for x in range(count):
        #     self.view.insert(edit, at, text)
        state = VintageState(self.view)
        if state.mode not in (MODE_VISUAL, MODE_VISUAL_LINE):
            for x in range(count):
                self.view.insert(edit, at, text)
        else:
            if text.endswith('\n'):
                text = text * count
                if not text.startswith('\n'):
                    text = '\n' + text
            else:
                text = text * count
            self.view.replace(edit, sel, text)


class ViEnterNormalMode(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)

        if state.mode == MODE_VISUAL:
            state.store_visual_selections()

        self.view.run_command('collapse_to_direction')
        self.view.run_command('dont_stay_on_eol_backward')
        state.enter_normal_mode()


class ViEnterNormalModeFromInsertMode(sublime_plugin.TextCommand):
    def run(self, edit):
        sels = list(self.view.sel())
        self.view.sel().clear()

        new_sels = []
        for s in sels:
            if s.a <= s.b:
                if (self.view.line(s.a).a != s.a):
                    new_sels.append(sublime.Region(s.a - 1, s.a - 1))
                else:
                    new_sels.append(sublime.Region(s.a, s.a))
            else:
                new_sels.append(s)

        for s in new_sels:
            self.view.sel().add(s)

        state = VintageState(self.view)
        state.enter_normal_mode()
        self.view.window().run_command('hide_auto_complete')


class ViEnterInsertMode(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)
        state.enter_insert_mode()
        self.view.run_command('collapse_to_direction')


class ViEnterVisualMode(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)
        state.enter_visual_mode()
        self.view.run_command('extend_to_minimal_width')


class ViEnterVisualLineMode(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)
        state.enter_visual_line_mode()


class ViEnterReplaceMode(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)
        state.enter_replace_mode()
        self.view.run_command('collapse_to_direction')
        state.reset()


class SetAction(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self, action):
        state = VintageState(self.view)
        state.action = action
        state.eval()


class SetMotion(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self, motion):
        state = VintageState(self.view)
        state.motion = motion
        state.eval()


class ViPushDigit(sublime_plugin.TextCommand):
    def run(self, edit, digit):
        state = VintageState(self.view)
        if not (state.action or state.motion):
            state.push_motion_digit(digit)
        elif state.action:
            state.push_action_digit(digit)


class ViReverseCaret(sublime_plugin.TextCommand):
    def run(self, edit):
        sels = list(self.view.sel())
        self.view.sel().clear()

        new_sels = []
        for s in sels:
            new_sels.append(sublime.Region(s.b, s.a))

        for s in new_sels:
            self.view.sel().add(s)


class ViEnterNormalInsertMode(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)
        state.enter_normal_insert_mode()

        # FIXME: We can't repeat 5ifoo<esc>
        self.view.run_command('mark_undo_groups_for_gluing')
        # ...User types text...


class ViRunNormalInsertModeActions(sublime_plugin.TextCommand):
    def run(self, edit):
        state = VintageState(self.view)
        # We've recorded what the user has typed into the buffer. Turn macro recording off.
        self.view.run_command('glue_marked_undo_groups')

        # FIXME: We can't repeat 5ifoo<esc> after we're done.
        for i in range(state.count - 1):
            self.view.run_command('repeat')

        # Ensure the count will be deleted.
        state.mode = MODE_NORMAL
        # Delete the count now.
        state.reset()

        self.view.run_command('vi_enter_normal_mode_from_insert_mode')


class SetRegister(sublime_plugin.TextCommand):
    def run(self, edit, character=None):
        state = VintageState(self.view)
        if character is None:
            state.expecting_register = True
        else:
            if character not in (REG_EXPRESSION,):
                state.register = character
                state.expecting_register = False
            else:
                self.view.run_command('vi_expression_register')


class ViExpressionRegister(sublime_plugin.TextCommand):
    def run(self, edit, insert=False, next_mode=None):
        def on_done(s):
            state = VintageState(self.view)
            try:
                rv = [str(eval(s, None, None)),]
                if not insert:
                    # TODO: We need to sort out the values received and sent to registers. When pasting,
                    # we assume a list... This should be encapsulated in Registers.
                    state.registers[REG_EXPRESSION] = rv
                else:
                    self.view.run_command('insert_snippet', {'contents': str(rv[0])})
                    state.reset()
            except:
                sublime.status_message("Vintageous: Invalid expression.")
                on_cancel()

        def on_cancel():
            state = VintageState(self.view)
            state.reset()

        self.view.window().show_input_panel('', '', on_done, None, on_cancel)


class ViR(sublime_plugin.TextCommand):
    def run(self, edit, character=None):
        state = VintageState(self.view)
        if character is None:
            state.action = 'vi_r'
            state.expecting_user_input = True
        else:
            state.user_input = character
            state.expecting_user_input= False
            state.eval()


class ViM(sublime_plugin.TextCommand):
    def run(self, edit, character=None):
        state = VintageState(self.view)
        state.action = 'vi_m'
        state.expecting_user_input = True


class _vi_m(sublime_plugin.TextCommand):
    def run(self, edit, character=None):
        state = VintageState(self.view)
        state.marks.add(character, self.view)


class ViQuote(sublime_plugin.TextCommand):
    def run(self, edit, character=None):
        state = VintageState(self.view)
        state.motion = 'vi_quote'
        state.expecting_user_input = True


class _vi_quote(sublime_plugin.TextCommand):
    def run(self, edit, mode=None, character=None, extend=False):
        def f(view, s):
            if mode == MODE_VISUAL:
                if s.a <= s.b:
                    if address.b < s.b:
                        return sublime.Region(s.a + 1, address.b)
                    else:
                        return sublime.Region(s.a, address.b)
                else:
                    return sublime.Region(s.a + 1, address.b)
            elif mode == MODE_NORMAL:
                return address
            elif mode == _MODE_INTERNAL_NORMAL:
                return sublime.Region(s.a, address.b)

            return s

        state = VintageState(self.view)
        address = state.marks.get_as_encoded_address(character)

        if address is None:
            return

        if isinstance(address, str):
            if not address.startswith('<command'):
                self.view.window().open_file(address, sublime.ENCODED_POSITION)
            else:
                # We get a command in this form: <command _vi_double_quote>
                self.view.run_command(address.split(' ')[1][:-1])
            return

        # This is a motion in a composite command.
        regions_transformer(self.view, f)


class ViF(sublime_plugin.TextCommand):
    def run(self, edit, character=None):
        state = VintageState(self.view)
        if character is None:
            state.motion = 'vi_f'
            state.expecting_user_input = True
        else:
            # FIXME: Dead code?
            state.user_input = character
            state.expecting_user_input= False
            state.eval()


class ViT(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    # XXX: Compare to ViBigF.
    def run(self, character=None):
        state = VintageState(self.view)
        if character is None:
            state.motion = 'vi_t'
            state.expecting_user_input = True
        else:
            state.user_input = character
            state.expecting_user_input= False
            state.eval()


class ViBigT(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    # XXX: Compare to ViBigF.
    def run(self, character=None):
        state = VintageState(self.view)
        if character is None:
            state.motion = 'vi_big_t'
            state.expecting_user_input = True
        else:
            state.user_input = character
            state.expecting_user_input= False
            state.eval()


class ViBigF(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self):
        state = VintageState(self.view)
        state.motion = 'vi_big_f'
        state.expecting_user_input = True


class ViI(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self, inclusive=False):
        state = VintageState(self.view)
        if inclusive:
            state.motion = 'vi_inclusive_text_object'
        else:
            state.motion = 'vi_exclusive_text_object'
        state.expecting_user_input = True


class CollectUserInput(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self, character=None):
        state = VintageState(self.view)
        state.user_input = character
        state.expecting_user_input= False
        state.eval()


class _vi_z_enter(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self):
        first_sel = self.view.sel()[0]
        current_row = self.view.rowcol(first_sel.b)[0] - 1

        topmost_visible_row, _ = self.view.rowcol(self.view.visible_region().a)

        self.view.run_command('scroll_lines', {'amount': (topmost_visible_row - current_row)})


class _vi_z_minus(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self):
        first_sel = self.view.sel()[0]
        current_row = self.view.rowcol(first_sel.b)[0]

        bottommost_visible_row, _ = self.view.rowcol(self.view.visible_region().b)

        number_of_lines = (bottommost_visible_row - current_row) - 1

        if number_of_lines > 1:
            self.view.run_command('scroll_lines', {'amount': number_of_lines})


class _vi_zz(IrreversibleTextCommand):
    def __init__(self, view):
        IrreversibleTextCommand.__init__(self, view)

    def run(self):
        first_sel = self.view.sel()[0]
        current_row = self.view.rowcol(first_sel.b)[0]

        topmost_visible_row, _ = self.view.rowcol(self.view.visible_region().a)
        bottommost_visible_row, _ = self.view.rowcol(self.view.visible_region().b)

        middle_row = (topmost_visible_row + bottommost_visible_row) / 2

        self.view.run_command('scroll_lines', {'amount': (middle_row - current_row)})


class _vi_r(sublime_plugin.TextCommand):
    def run(self, edit, character=None, mode=None):
        if mode == _MODE_INTERNAL_NORMAL:
            for s in self.view.sel():
                self.view.replace(edit, s, character * s.size())


class _vi_undo(IrreversibleTextCommand):
    """Once the latest vi command has been undone, we might be left with non-empty selections.
    This is due to the fact that Vintageous defines selections in a separate step to the actual
    command running. For example, v,e,d,u would undo the deletion operation and restore the
    selection that v,e had created.

    Assuming that after an undo we're back in normal mode, we can take for granted that any leftover
    selections must be destroyed. I cannot think of any situation where Vim would have to restore
    selections after *u*, but it may well happen under certain circumstances I'm not aware of.

    Note 1: We are also relying on Sublime Text to restore the v or V selections existing at the
    time the edit command was run. This seems to be safe, but we're blindly relying on it.

    Note 2: Vim knows the position the caret was in before creating the visual selection. In
    Sublime Text we lose that information (at least it doesn't seem to be straightforward to
    obtain).
    """
    #  !!! This is a special command that does not go through the usual processing. !!!
    #  !!! It must skip the undo stack. !!!

    # TODO: It must be possible store or retrieve the actual position of the caret before the
    # visual selection performed by the user.
    def run(self):
        # We define our own transformer here because we want to handle undo as a special case.
        # TODO: I don't know if it needs to be an special case in reality.
        def f(view, s):
            # Compensates the move issued below.
            if s.a < s.b :
                return sublime.Region(s.a + 1, s.a + 1)
            else:
                return sublime.Region(s.a, s.a)

        state = VintageState(self.view)
        for i in range(state.count):
            self.view.run_command('undo')

        if self.view.has_non_empty_selection_region():
            regions_transformer(self.view, f)
            # !! HACK !! /////////////////////////////////////////////////////////
            # This is a hack to work around an issue in Sublime Text:
            # When undoing in normal mode, Sublime Text seems to prime a move by chars
            # forward that has never been requested by the user or Vintageous.
            # As far as I can tell, Vintageous isn't at fault here, but it seems weird
            # to think that Sublime Text is wrong.
            self.view.run_command('move', {'by': 'characters', 'forward': False})
            # ////////////////////////////////////////////////////////////////////

        state.update_xpos()
        # Ensure that we wipe the count, if any.
        state.reset()


class _vi_repeat(IrreversibleTextCommand):
    """Vintageous manages the repeat operation on its own to ensure that we always use the latest
       modifying command, instead of being tied to the undo stack (as Sublime Text is by default).
    """

    #  !!! This is a special command that does not go through the usual processing. !!!
    #  !!! It must skip the undo stack. !!!

    def run(self):
        state = VintageState(self.view)

        try:
            cmd, args, _ = state.repeat_command
        except TypeError:
            # Unreachable.
            return

        if not cmd:
            return
        elif cmd == 'vi_run':
            args['next_mode'] = MODE_NORMAL
            args['follow_up_mode'] = 'vi_enter_normal_mode'
            args['count'] = state.count * args['count']
            self.view.run_command(cmd, args)
        elif cmd == 'sequence':
            for i, _ in enumerate(args['commands']):
                # Access this shape: {"commands":[['vi_run', {"foo": 100}],...]}
                args['commands'][i][1]['next_mode'] = MODE_NORMAL
                args['commands'][i][1]['follow_up_mode'] = 'vi_enter_normal_mode'

            # TODO: Implement counts properly for 'sequence' command.
            for i in range(state.count):
                self.view.run_command(cmd, args)

        # Ensure we wipe count data if any.
        state.reset()
        # XXX: Needed here? Maybe enter_... type commands should be IrreversibleCommands so we
        # must/can call them whenever we need them withouth affecting the undo stack.
        self.view.run_command('vi_enter_normal_mode')


class _vi_ctrl_w_v_action(sublime_plugin.TextCommand):
    def run(self, edit):
        self.view.window().run_command('new_pane', {})


class Sequence(sublime_plugin.TextCommand):
    """Required so that mark_undo_groups_for_gluing and friends work.
    """
    def run(self, edit, commands):
        for cmd, args in commands:
            self.view.run_command(cmd, args)

        # XXX: Sequence is a special case in that it doesn't run through vi_run, so we need to
        # ensure the next mode is correct. Maybe we can improve this by making it more similar to
        # regular commands?
        state = VintageState(self.view)
        state.enter_normal_mode()


class _vi_big_j(sublime_plugin.TextCommand):
    def run(self, edit, mode=None):
        def f(view, s):
            if mode == _MODE_INTERNAL_NORMAL:
                full_current_line = view.full_line(s.b)
                target = full_current_line.b - 1
                full_next_line = view.full_line(full_current_line.b)
                two_lines = sublime.Region(full_current_line.a, full_next_line.b)

                # Text without \n.
                first_line_text = view.substr(view.line(full_current_line.a))
                next_line_text = view.substr(full_next_line)

                if len(next_line_text) > 1:
                    next_line_text = next_line_text.lstrip()

                sep = ''
                if first_line_text and not first_line_text.endswith(' '):
                    sep = ' '

                view.replace(edit, two_lines, first_line_text + sep + next_line_text)

                if first_line_text:
                    return sublime.Region(target, target)
                return s
            else:
                return s

        regions_transformer(self.view, f)


class _vi_ctrl_a(sublime_plugin.TextCommand):
    def run(self, edit, count=1, mode=None):
        def f(view, s):
            if mode == _MODE_INTERNAL_NORMAL:
                word = view.word(s.a)
                new_digit = int(view.substr(word)) + count
                view.replace(edit, word, str(new_digit))

            return s

        if mode != _MODE_INTERNAL_NORMAL:
            return

        # TODO: Deal with octal, hex notations.
        # TODO: Improve detection of numbers.
        # TODO: Find the next numeric word in the line if none is found under the caret.
        words = [self.view.substr(self.view.word(s)) for s in self.view.sel()]
        if not all([w.isdigit() for w in words]):
            utils.blink()
            return

        regions_transformer(self.view, f)


class _vi_ctrl_x(sublime_plugin.TextCommand):
    def run(self, edit, count=1, mode=None):
        def f(view, s):
            if mode == _MODE_INTERNAL_NORMAL:
                word = view.word(s.a)
                new_digit = int(view.substr(word)) - count
                view.replace(edit, word, str(new_digit))

            return s

        if mode != _MODE_INTERNAL_NORMAL:
            return

        # TODO: Deal with octal, hex notations.
        # TODO: Improve detection of numbers.
        # TODO: Find the next numeric word in the line if none is found under the caret.
        words = [self.view.substr(self.view.word(s)) for s in self.view.sel()]
        if not all([w.isdigit() for w in words]):
            utils.blink()
            return

        regions_transformer(self.view, f)


class _vi_g_v(IrreversibleTextCommand):
    def run(self):
        # Assume normal mode.
        regs = (self.view.get_regions('vi_visual_selections') or
                list(self.view.sel()))

        self.view.sel().clear()
        for r in regs:
            self.view.sel().add(r)


class ViQ(IrreversibleTextCommand):
    def run(self):
        state = VintageState(self.view)
        state.action = 'vi_q'
        state.expecting_user_input = True


class _vi_q(IrreversibleTextCommand):
    def run(self, name=None):
        state = VintageState(self.view)

        if name == None and not state.is_recording:
            return

        if not state.is_recording:
            state._latest_macro_name = name
            state.is_recording = True
            self.view.run_command('start_record_macro')
            return

        if state.is_recording:
            self.view.run_command('stop_record_macro')
            state.is_recording = False
            state.reset()

            # Store the macro away.
            modifying_cmd = self.view.command_history(0, True)
            state.latest_macro = modifying_cmd


class _vi_run_macro(IrreversibleTextCommand):
    def run(self, name=None):
        if not (name and VintageState(self.view).latest_macro):
            return

        if name == '@':
            # Run the macro recorded latest.
            self.view.run_command('run_macro')
        else:
            # TODO: Implement macro registers.
            self.view.run_command('run_command')


class ViAt(IrreversibleTextCommand):
    def run(self):
        state = VintageState(self.view)
        state.action = 'vi_at'
        state.expecting_user_input = True
