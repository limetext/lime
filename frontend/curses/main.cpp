#include <ncurses.h>
#include <string.h>
#include <locale.h>

#include <ApplicationServices/ApplicationServices.h>
#include <Carbon/Carbon.h>
#include <CoreServices/CoreServices.h>
#include <pthread.h>

class EventTapKeyControl
{
public:
    EventTapKeyControl()
    : hasInput(true), running(true)
    {
        int id = getTerminalPid();
        ProcessSerialNumber currentProcess = {kNoProcess, kNoProcess};

        while (true)
        {
            GetNextProcess(&currentProcess);
            CFDictionaryRef dict = ProcessInformationCopyDictionary(&currentProcess, kProcessDictionaryIncludeAllInformationMask);

            if (dict)
            {
                //CFDictionaryApplyFunction(dict, bice, NULL);
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
private:

    pthread_t thread;
    bool hasInput;
    bool running;

    static void* mainloop(void* data)
    {
        EventTapKeyControl* ctl = (EventTapKeyControl*) data;
        while (ctl->running)
        {
            int ch = getch();
            ctl->hasInput = true;
        }
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

    static CGEventRef myCGEventCallback(CGEventTapProxy proxy, CGEventType type,  CGEventRef event, void* refcon)
    {
        EventTapKeyControl* ctl = (EventTapKeyControl*) refcon;

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

        if (flags & kCGEventFlagMaskControl)
        {
            strcat(keys, "ctrl ");
        }

        if (flags & kCGEventFlagMaskAlternate)
        {
            strcat(keys, "alt ");
        }

        if (flags & kCGEventFlagMaskShift)
        {
            strcat(keys, "shift ");
        }

        if (flags & kCGEventFlagMaskCommand)
        {
            strcat(keys, "cmd ");
        }

        UInt32 code1 = ctl->getCharForKey(rawcode);
        UInt32 code2 = ctl->getCharForKey(rawcode, flags);

        if (code1 == 'c' && (flags & kCGEventFlagMaskControl))
        {
            ctl->running = false;
            CFRunLoopStop(CFRunLoopGetCurrent());
        }

        char buf[512];

        if (code1 == 0x1b)
        {
            sprintf(buf, "%s escape", keys);
        }
        else if (code1 == 0x1c)
        {
            sprintf(buf, "%s left", keys);
        }
        else if (code1 == 0x1d)
        {
            sprintf(buf, "%s right", keys);
        }
        else if (code1 == 0x1e)
        {
            sprintf(buf, "%s up", keys);
        }
        else if (code1 == 0x1f)
        {
            sprintf(buf, "%s down", keys);
        }
        else
        {
            sprintf(buf, "%s %lc (%d) %lc (%d)", keys, (wchar_t) code1, code1, (wchar_t) code2, code2);
        }

        int row, col;
        getmaxyx(stdscr, row, col);
        clear();
        mvprintw(row / 2, (col - strlen(buf)) / 2, buf, "");
        refresh();

        return NULL;
    }
};
int main(int argc, const char* argv[])
{

    setlocale(LC_ALL, "");

    initscr();
    raw();
    noecho();
    mousemask(ALL_MOUSE_EVENTS, NULL);
    keypad(stdscr, true);

    EventTapKeyControl c;

    endwin();
    printf("cleaning up\n");

    return 0;
}
