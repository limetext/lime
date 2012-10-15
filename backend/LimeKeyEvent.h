#ifndef __INCLUDED_LIME_KEYEVENT_H
#define __INCLUDED_LIME_KEYEVENT_H

enum
{
    LIME_KEY_MODIFIER_NONE    = 0,
    LIME_KEY_MODIFIER_SHIFT   = (1<<0),
    LIME_KEY_MODIFIER_ALT     = (1<<1),
    LIME_KEY_MODIFIER_CTRL    = (1<<2),
    LIME_KEY_MODIFIER_COMMAND = (1<<3),
};
typedef int LimeKeyModifiers;

class LimeKeyEvent
{
public:
    LimeKeyEvent(wchar_t unicode=0, wchar_t raw=0, LimeKeyModifiers mod=LIME_KEY_MODIFIER_NONE)
    :   mUnicodeChar(unicode), mRawChar(raw), mModifiers(mod)
    {
    }

    wchar_t GetUnicodeChar() const
    {
        return mUnicodeChar;
    }
    void SetUnicodeChar(wchar_t arg)
    {
        mUnicodeChar = arg;
    }
    wchar_t GetRawChar() const
    {
        return mRawChar;
    }
    void SetRawChar(wchar_t arg)
    {
        mRawChar = arg;
    }
    LimeKeyModifiers GetModifiers() const
    {
        return mModifiers;
    }
    void SetModifiers(LimeKeyModifiers arg)
    {
        mModifiers = arg;
    }

private:

    LimeKeyModifiers mModifiers;
    wchar_t mUnicodeChar;
    wchar_t mRawChar;
};

#endif
