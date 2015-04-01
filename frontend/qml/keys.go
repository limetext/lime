// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package main

import (
	"github.com/limetext/lime/backend/keys"
)

// http://qt-project.org/doc/qt-5.1/qtcore/qt.html#Key-enum
var lut = map[int]keys.Key{
	0x01000000: keys.Escape,
	0x01000001: '\t',
	// 0x01000002 // Qt::Key_Backtab
	0x01000003: keys.Backspace,
	0x01000004: keys.Enter,
	0x01000005: keys.KeypadEnter,
	0x01000006: keys.Insert,
	0x01000007: keys.Delete,
	0x01000008: keys.Break,
	// 0x01000009 // Qt::Key_Print
	// 0x0100000a // Qt::Key_SysReq
	// 0x0100000b // Qt::Key_Clear
	0x01000010: keys.Home,
	0x01000011: keys.End,
	0x01000012: keys.Left,
	0x01000013: keys.Up,
	0x01000014: keys.Right,
	0x01000015: keys.Down,
	0x01000016: keys.PageUp,
	0x01000017: keys.PageDown,
	// 0x01000020 // Qt::Key_Shift
	// 0x01000021 // Qt::Key_Control On Mac OS X, this corresponds to the Command keys.
	// 0x01000022 // Qt::Key_Meta On Mac OS X, this corresponds to the Control keys. On Windows keyboards, this key is mapped to the Windows key.
	// 0x01000023 // Qt::Key_Alt
	// 0x01001103 // Qt::Key_AltGr On Windows, when the KeyDown event for this key is sent, the Ctrl+Alt modifiers are also set.
	// 0x01000024 // Qt::Key_CapsLock
	// 0x01000025 // Qt::Key_NumLock
	// 0x01000026 // Qt::Key_ScrollLock
	0x01000030: keys.F1,
	0x01000031: keys.F2,
	0x01000032: keys.F3,
	0x01000033: keys.F4,
	0x01000034: keys.F5,
	0x01000035: keys.F6,
	0x01000036: keys.F7,
	0x01000037: keys.F8,
	0x01000038: keys.F9,
	0x01000039: keys.F10,
	0x0100003a: keys.F11,
	0x0100003b: keys.F12,
	// 0x01000053 // Qt::Key_Super_L
	// 0x01000054 // Qt::Key_Super_R
	// 0x01000055 // Qt::Key_Menu
	// 0x01000056 // Qt::Key_Hyper_L
	// 0x01000057 // Qt::Key_Hyper_R
	// 0x01000058 // Qt::Key_Help
	// 0x01000059 // Qt::Key_Clear
	// 0x01000060 // Qt::Key_Direction_R
	0x20: ' ',
	0x21: '!',
	0x22: '"',
	0x23: '#',
	0x24: '$',
	0x25: '%',
	0x26: '&',
	0x27: '\'',
	0x28: '(',
	0x29: ')',
	0x2a: '*',
	0x2b: '+',
	0x2c: ',',
	0x2d: '-',
	0x2e: '.',
	0x2f: '/',
	0x30: '0',
	0x31: '1',
	0x32: '2',
	0x33: '3',
	0x34: '4',
	0x35: '5',
	0x36: '6',
	0x37: '7',
	0x38: '8',
	0x39: '9',
	0x3a: ':',
	0x3b: ';',
	0x3c: '<',
	0x3d: '=',
	0x3e: '>',
	0x3f: '?',
	0x40: '@',
	0x41: 'a',
	0x42: 'b',
	0x43: 'c',
	0x44: 'd',
	0x45: 'e',
	0x46: 'f',
	0x47: 'g',
	0x48: 'h',
	0x49: 'i',
	0x4A: 'j',
	0x4B: 'k',
	0x4C: 'l',
	0x4d: 'm',
	0x4e: 'n',
	0x4f: 'o',
	0x50: 'p',
	0x51: 'q',
	0x52: 'r',
	0x53: 's',
	0x54: 't',
	0x55: 'u',
	0x56: 'v',
	0x57: 'w',
	0x58: 'x',
	0x59: 'y',
	0x5a: 'z',
	0x5b: '[',
	0x5c: '\\',
	0x5d: ']',
	0x5e: '°', // Qt::Key_AsciiCircum
	0x5f: '_', // Qt::Key_Underscore
	0x60: '`', // Qt::Key_QuoteLeft
	0x7b: '{', // Qt::Key_BraceLeft
	0x7c: '|', // Qt::Key_Bar
	0x7d: '}', // Qt::Key_BraceRight
	0x7e: '~', // Qt::Key_AsciiTilde
	// 0x0a0: '', // Qt::Key_nobreakspace
	// 0x0a1: '', // Qt::Key_exclamdown
	// 0x0a2: '', // Qt::Key_cent
	// 0x0a3: '', // Qt::Key_sterling
	// 0x0a4: '', // Qt::Key_currency
	// 0x0a5: '', // Qt::Key_yen
	// 0x0a6: '', // Qt::Key_brokenbar
	// 0x0a7: '', // Qt::Key_section
	// 0x0a8: '', // Qt::Key_diaeresis
	// 0x0a9: '', // Qt::Key_copyright
	// 0x0aa: '', // Qt::Key_ordfeminine
	// 0x0ab: '', // Qt::Key_guillemotleft
	// 0x0ac: '', // Qt::Key_notsign
	// 0x0ad: '', // Qt::Key_hyphen
	// 0x0ae: '', // Qt::Key_registered
	// 0x0af: '', // Qt::Key_macron
	0x0b0: '°', // Qt::Key_degree
	// 0x0b1: '', // Qt::Key_plusminus
	0x0b2: '²', // Qt::Key_twosuperior
	0x0b3: '³', // Qt::Key_threesuperior
	0x0b4: '´', // Qt::Key_acute
	// 0x0b5: '', // Qt::Key_mu
	// 0x0b6: '', // Qt::Key_paragraph
	// 0x0b7: '', // Qt::Key_periodcentered
	// 0x0b8: '', // Qt::Key_cedilla
	// 0x0b9: '', // Qt::Key_onesuperior
	// 0x0ba: '', // Qt::Key_masculine
	// 0x0bb: '', // Qt::Key_guillemotright
	// 0x0bc: '', // Qt::Key_onequarter
	// 0x0bd: '', // Qt::Key_onehalf
	// 0x0be: '', // Qt::Key_threequarters
	// 0x0bf: '', // Qt::Key_questiondown
	// 0x0c0: '', // Qt::Key_Agrave
	// 0x0c1: '', // Qt::Key_Aacute
	// 0x0c2: '', // Qt::Key_Acircumflex
	// 0x0c3: '', // Qt::Key_Atilde
	0x0c4: 'ä', // Qt::Key_Adiaeresis
	// 0x0c5: '', // Qt::Key_Aring
	// 0x0c6: '', // Qt::Key_AE
	// 0x0c7: '', // Qt::Key_Ccedilla
	// 0x0c8: '', // Qt::Key_Egrave
	// 0x0c9: '', // Qt::Key_Eacute
	// 0x0ca: '', // Qt::Key_Ecircumflex
	// 0x0cb: '', // Qt::Key_Ediaeresis
	// 0x0cc: '', // Qt::Key_Igrave
	// 0x0cd: '', // Qt::Key_Iacute
	// 0x0ce: '', // Qt::Key_Icircumflex
	// 0x0cf: '', // Qt::Key_Idiaeresis
	// 0x0d0: '', // Qt::Key_ETH
	// 0x0d1: '', // Qt::Key_Ntilde
	// 0x0d2: '', // Qt::Key_Ograve
	// 0x0d3: '', // Qt::Key_Oacute
	// 0x0d4: '', // Qt::Key_Ocircumflex
	// 0x0d5: '', // Qt::Key_Otilde
	0x0d6: 'ö', // Qt::Key_Odiaeresis
	// 0x0d7: '', // Qt::Key_multiply
	// 0x0d8: '', // Qt::Key_Ooblique
	// 0x0d9: '', // Qt::Key_Ugrave
	// 0x0da: '', // Qt::Key_Uacute
	// 0x0db: '', // Qt::Key_Ucircumflex
	0x0dc: 'ü', // Qt::Key_Udiaeresis
	// 0x0dd: '', // Qt::Key_Yacute
	// 0x0de: '', // Qt::Key_THORN
	// 0x0df: '', // Qt::Key_ssharp
	// 0x0f7: '', // Qt::Key_division
}
