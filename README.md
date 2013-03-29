# What is Lime?

I love the [Sublime Text](http://www.sublimetext.com) editor. [I have](https://github.com/quarnster/SublimeClang) [created](https://github.com/quarnster/SublimeJava) [several](https://github.com/quarnster/CompleteSharp) [plugins](https://github.com/quarnster/SublimeGDB) [to make](https://github.com/quarnster/ADBView) it even better. One thing that scares me though is that it is not open sourced and the [pace of nightly releases](http://www.sublimetext.com/nightly) have recently been anything but nightly.

So I started thinking about what it would take to create an open source clone from scratch and how I would go about doing it. After browsing around a bit for high quality font rendering engines a thought hit me; why not just use a web browser engine such as [WebKit](http://www.webkit.org/) as the backend? A Web browser engine already solves a lot of problems when it comes to layout, coloring and high quality font rendering and is already up and running on many platforms so it takes care of making the editor run everywhere too.

Then I started thinking, "well, what if it ran in the terminal?" so I started writing some code for that too.

At this point I don't know what direction it'll be heading and how serious it'll become. I consider it very experimental with possible multiple frontends (terminal, gui app, html) initially served by a Sublime Text compatible backend, but later on it'll fix some of the annoyances I have with Sublime Text.
