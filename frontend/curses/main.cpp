#include <ncurses.h>
#include <string.h>
#include <locale.h>

#include <Carbon/Carbon.h>
#include <pthread.h>

#include "LimeKeyEvent.h"
#include "Syntax.h"

#ifdef HAVE_TERMKEY
#include <termkey.h>
#endif

class EventTapKeyControl
{
public:
    EventTapKeyControl()
    : hasInput(true), running(true), flags(0)
    {
        int id = getTerminalPid();
        ProcessSerialNumber currentProcess = {kNoProcess, kNoProcess};

        while (true)
        {
            GetNextProcess(&currentProcess);
            CFDictionaryRef dict = ProcessInformationCopyDictionary(&currentProcess, kProcessDictionaryIncludeAllInformationMask);

            if (dict)
            {
                CFNumberRef ref = (CFNumberRef) CFDictionaryGetValue(dict, CFSTR("pid"));

                if (ref)
                {
                    int id2;
                    CFNumberGetValue(ref, kCFNumberIntType, &id2);

                    if (id == id2)
                    {
                        break;
                    }
                }

                CFRelease(dict);
            }
            else
            {
                throw "Couldn't get terminal process serial number";
            }
        }
#ifdef HAVE_TERMKEY
        tk = termkey_new(0, 0);
#endif

        pthread_create(&thread, NULL, mainloop, this);

        CFMachPortRef      eventTap;
        CFRunLoopSourceRef runLoopSource;
        CGEventMask        eventMask = kCGEventMaskForAllEvents;

        eventTap = CGEventTapCreateForPSN(&currentProcess, kCGHeadInsertEventTap, 0, eventMask, myCGEventCallback, this);
        runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
        CFRunLoopAddSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
        CGEventTapEnable(eventTap, true);
        CFRunLoopRun();
    }

    ~EventTapKeyControl()
    {
        running = false;
        pthread_join(thread, NULL);
#ifdef HAVE_TERMKEY
        termkey_destroy(tk);
#endif
    }

private:

#ifdef HAVE_TERMKEY
    TermKey *tk;
#endif

    LimeKeyEvent currentEvent;
    pthread_t thread;
    bool hasInput;
    bool running;
    CGEventFlags flags;

    static void* mainloop(void* data)
    {
        EventTapKeyControl* ctl = (EventTapKeyControl*) data;
        while (ctl->running)
        {
            char buf[512];
            LimeKeyModifiers mod = 0;
            bool clearMod = true;

#ifdef HAVE_TERMKEY
            TermKeyResult ret;
            TermKeyKey key;

            ret = termkey_waitkey(ctl->tk, &key);
            termkey_strfkey(ctl->tk, buf, sizeof buf, &key, TERMKEY_FORMAT_LONGMOD);
            ctl->currentEvent.SetRawChar(tolower(key.code.number));
            ctl->currentEvent.SetUnicodeChar(key.code.number);

            mod = ctl->currentEvent.GetModifiers();
            clearMod = mod == 0;

            if (clearMod)
            {
                if (key.modifiers & TERMKEY_KEYMOD_CTRL)
                    mod |= LIME_KEY_MODIFIER_CTRL;
                if (key.modifiers & TERMKEY_KEYMOD_ALT)
                    mod |= LIME_KEY_MODIFIER_ALT;
                if (key.modifiers & TERMKEY_KEYMOD_SHIFT)
                    mod |= LIME_KEY_MODIFIER_SHIFT;
            }
            ctl->currentEvent.SetModifiers(mod);
#else
            int ch = getch();
            ctl->hasInput = true;
            mod = ctl->currentEvent.GetModifiers();
            clearMod = mod == 0;
            switch (ch)
            {
                case 0x3: // ctrl+c
                    ch = 'c';
                    mod |= LIME_KEY_MODIFIER_CTRL;
                    break;
                default:
                    break;
            }
            ctl->currentEvent.SetModifiers(mod);

            ctl->currentEvent.SetUnicodeChar(ch);
            ctl->currentEvent.SetRawChar(ch);

#endif
            // TODO: dispatch
            sprintf(buf, "%d %d %lc (%d) %lc (%d)", ctl->currentEvent.GetModifiers(),
               ctl->currentEvent.GetRawChar(), ctl->currentEvent.GetRawChar(),
               ctl->currentEvent.GetUnicodeChar(), ctl->currentEvent.GetUnicodeChar()
            );

            if (ctl->currentEvent.GetModifiers() == LIME_KEY_MODIFIER_CTRL &&
                ctl->currentEvent.GetRawChar() == 'c')
            {
                ctl->running = false;
            }

            int row, col;
            getmaxyx(stdscr, row, col);
            clear();
            mvprintw(row / 2, (col - strlen(buf)) / 2, buf, "");
            refresh();
            ctl->currentEvent = LimeKeyEvent(0, 0, clearMod ? 0 : mod);
        }
        return NULL;
    }

    int getTerminalPid()
    {
        int curr = getppid();

        while (true)
        {
            char buf[512];
            sprintf(buf, "ps -p %d -o ppid,command | tail +2", curr);
            FILE* fp = popen(buf, "r");

            if (fp)
            {
                if (fgets(buf, 512, fp) != NULL)
                {
                    int ppid;
                    char app[512];
                    sscanf(buf, "%d %s", &ppid, app);

                    if (strstr(app, "/Applications/"))
                    {
                        return curr;
                    }

                    curr = ppid;
                }
                else
                {
                    break;
                }

                pclose(fp);
            }
        }

        return -1;
    }

    UInt32 getCharForKey(CGKeyCode keyCode, UInt32 modifierFlags = 0)
    {
        TISInputSourceRef currentKeyboard = TISCopyCurrentKeyboardInputSource();
        CFDataRef layoutData = (CFDataRef) TISGetInputSourceProperty(currentKeyboard, kTISPropertyUnicodeKeyLayoutData);
        const UCKeyboardLayout* keyboardLayout = (const UCKeyboardLayout*)CFDataGetBytePtr(layoutData);

        UInt32 keysDown = 0;
        UniChar chars[4];
        UniCharCount realLength;

        UCKeyTranslate(keyboardLayout,
                       keyCode,
                       kUCKeyActionDisplay,
                       modifierFlags,
                       LMGetKbdType(),
                       kUCKeyTranslateNoDeadKeysBit,
                       &keysDown,
                       sizeof(chars) / sizeof(chars[0]),
                       &realLength,
                       chars);
        CFRelease(currentKeyboard);

        if (realLength == 1)
        {
            return chars[0];
        }
        else
        {
            return 0;
        }
    }

    LimeKeyModifiers GetModifiers(CGEventFlags flags)
    {
        LimeKeyModifiers mod = LIME_KEY_MODIFIER_NONE;
        if (flags & kCGEventFlagMaskAlternate)
            mod |= LIME_KEY_MODIFIER_ALT;
        if (flags & kCGEventFlagMaskShift)
            mod |= LIME_KEY_MODIFIER_SHIFT;
        if (flags & kCGEventFlagMaskCommand)
            mod |= LIME_KEY_MODIFIER_COMMAND;
        if (flags & kCGEventFlagMaskControl)
            mod |= LIME_KEY_MODIFIER_CTRL;
        return mod;
    }

    static CGEventRef myCGEventCallback(CGEventTapProxy proxy, CGEventType type,  CGEventRef event, void* refcon)
    {
        EventTapKeyControl* ctl = (EventTapKeyControl*) refcon;
        ctl->flags = CGEventGetFlags(event);

        if (type == kCGEventFlagsChanged)
        {
            mvprintw(0, 0, "mod: %d\t\t\t\t", ctl->GetModifiers(ctl->flags));
            refresh();
            ctl->currentEvent.SetModifiers(ctl->GetModifiers(ctl->flags));
        }

        if (!ctl->running)
        {
            CFRunLoopStop(CFRunLoopGetCurrent());
        }

        if (!ctl->hasInput)
            return event;

        switch (type)
        {
            case kCGEventKeyDown:
            case kCGEventKeyUp:
                break;
            case kCGEventLeftMouseDown:
            case kCGEventLeftMouseUp:
            case kCGEventMouseMoved:
                ctl->hasInput = false;
                return event;
            default:
                return event;
        }

        CGKeyCode rawcode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
        CGEventFlags flags = CGEventGetFlags(event);
        char keys[512] = "";

        ctl->currentEvent.SetUnicodeChar(ctl->getCharForKey(rawcode, flags));
        ctl->currentEvent.SetRawChar(ctl->getCharForKey(rawcode));
        ctl->currentEvent.SetModifiers(ctl->GetModifiers(ctl->flags));
        // TODO: dispatch.
        // return NULL;

        return event;
    }
};
int main(int argc, const char* argv[])
{
    lime::backend::Syntax s("../3rdparty/javascript.tmbundle/Syntaxes/JavaScript.plist");

    setlocale(LC_ALL, "");

    initscr();
#ifndef HAVE_TERMKEY
    raw();
    noecho();
    mousemask(ALL_MOUSE_EVENTS, NULL);
    keypad(stdscr, true);
#endif

    {
        EventTapKeyControl c;
    }
#ifndef HAVE_TERMKEY
    noraw();
    echo();
    keypad(stdscr, false);
#endif
    endwin();
    printf("cleaning up\n");

    return 0;
}
