# What is Lime?

I love the [Sublime Text 2](http://www.sublimetext.com) editor. [I have](https://github.com/quarnster/SublimeClang) [created](https://github.com/quarnster/SublimeJava) [several](https://github.com/quarnster/CompleteSharp) [plugins](https://github.com/quarnster/SublimeGDB) [to make](https://github.com/quarnster/ADBView) it even better. One thing that scares me though is that it is not open sourced and the [pace of nightly releases](http://www.sublimetext.com/nightly) have recently been anything but nightly.

So I started thinking about what it would take to create an open source clone from scratch and how I would go about doing it. After browsing around a bit for high quality font rendering engines a thought hit me; why not just use a web browser engine such as [WebKit](http://www.webkit.org/) as the backend? A Web browser engine already solves a lot of problems when it comes to layout, coloring and high quality font rendering and is already up and running on many platforms so it takes care of making the editor run everywhere too.

It didn't take long before another thought hit me; what if it rather than using a web browser engine as the backend, it ran directly inside your regular browser? The more I thought about it, the more I liked the idea and thus Lime was born.

