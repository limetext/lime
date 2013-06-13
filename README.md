# What is Lime?

I love the [Sublime Text](http://www.sublimetext.com) editor. [I have](https://github.com/quarnster/SublimeClang) [created](https://github.com/quarnster/SublimeJava) [several](https://github.com/quarnster/CompleteSharp) [plugins](https://github.com/quarnster/SublimeGDB) [to make](https://github.com/quarnster/ADBView) it even better. One thing that scares me though is that it is not open sourced and the [pace of nightly releases](http://www.sublimetext.com/nightly) have recently been anything but nightly, even now that version 3 is out in Beta.

There was a period of about 6 months after the Sublime Text 2 "stable" version was released where pretty much nothing at all was communicated to the users about what to expect in the future, nor was there much support offered in the forums. People including myself were wondering if the product was dead and I personally wondered what would happen to all the bugs, crashes and annoyances that still existed in ST2. This lack of communication is a dealbreaker to me and I decided that I will not spend any more money on that product because of it.

As none of the other text editors I've tried come close to the love I had for Sublime Text, I decided I had to create my own.

# Goals

- [x] 100% Open source
- [x] Compatible with Textmate color schemes (which is what ST is using)
- [x] Compatible with Textmate syntax definitions (which again is what ST is using)
- [ ] Compatible with Textmate snippets
- [ ] Compatible with Sublime Text's python plugin api. I'll probably never implement this 100%, only the api bits I need for the plugins I use.
- [ ] Compatible with Sublime Text's keybindings and settings
- [ ] Compatible with Sublime Text snippets
- [ ] Sublime Text's Goto anything panel
- [x] Multiple cursors
- [x] Regression tests (Programming in [Go](http://golang.org) makes it trivial and even fun to write them ;))
- [ ] Support for plugging in a custom parser for more advanced syntax highlighting.
- [ ] Terminal UI (*Maybe* I'll work on a simple non-terminal UI at some point)
- [ ] Cross platform (It appears to be compiling and running on OSX and Linux last I tried, but needs further validation.)

# Why can't I open up an issue?

Because I'm just a single person and I don't want to offer up my spare time doing support or dealing with feature requests that I don't care about myself. If you want a feature implemented or a bug fixed, fork it and implement it yourself and submit a pull request when you're happy with the implementation.

# Build instructions

You need to have Go 1.1 installed. As of writing this it hasn't been released yet which means you need to [build it yourself](http://tip.golang.org/doc/install/source). If you haven't built it already, you might be interested to know that on my machine it takes less than 60 seconds to build the whole Go distribution including the compilers and all packages in the standard library, so this is not a time consuming task. Go is the absolutely easiest and fastest system language to get a full environment up and running for from source that I've ever stumbled upon, without huge dependencies on outside resources. IIRC all you need is a c compiler installed.

Once go is installed and set up properly a rough draft is (please submit a pull request if you find other steps are needed):

```
go get code.google.com/p/log4go github.com/quarnster/parser github.com/quarnster/completion
sudo apt-get install libonig-dev python3-dev (on Linux)
brew install oniguruma python3 (on OSX)
git clone --recursive git@github.com:quarnster/lime.git
cd lime/3rdparty/libs/gopy/lib (Tweak cgo.go as appropriate with the help of python3-config --cflags and python3-config --libs)
cd lime/build
go run build.go
cd ../frontend/termbox
go run main.go
```

# License

The license of the project is the 2-clause BSD license:

```
Copyright (c) 2013 Fredrik Ehnbom
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
```
