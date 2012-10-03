/*
Copyright 2012 Fredrik Ehnbom

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

function Theme(name)
{
    var tmLang = loadFile(name);
    this.jsonString = PlistParser.parse(toXML(tmLang));
    var cssDef = "";

    this.createCss = function(name, setting)
    {
        cssDef += name + "\n{\n";

        if (setting.settings.foreground)
        {
            cssDef += "\tcolor:" + setting.settings.foreground + ";\n";
        }
        if (setting.settings.background)
        {
            cssDef += "\tbackground-color:" + setting.settings.background + ";\n";
        }

        cssDef += "}\n";
    }

    for (var i in this.jsonString.settings)
    {
        var setting = this.jsonString.settings[i];
        if (setting.settings)
        {
            var name = "body";
            if (setting.scope)
            {
                setting.scope = setting.scope.split(",");
                for (var j in setting.scope)
                {
                    setting.scope[j] = setting.scope[j].trim();
                    name = "." + setting.scope[j].replace(/\./g, "_");
                    this.createCss(name, setting);
                }
            }
            else
            {
                this.createCss(name, setting);
            }
        }
        cssDef += ".default\n{\n";
        cssDef += "\tfont-family:\"Menlo Regular\", monospace;\n";
        cssDef += "\tfont-size:12px;\n";
        cssDef += "}";
    }
    var sheet = document.createElement('style')
    sheet.innerHTML = cssDef;
    document.body.appendChild(sheet);


    this.getCssClassesForScopes = function(scopes)
    {
        while (scopes.length)
        {
            for (var i in this.jsonString.settings)
            {
                var setting = this.jsonString.settings[i];
                if (setting.scope)
                {
                    for (var j in setting.scope)
                    {
                        if (scopes.endsWith(setting.scope[j]))
                            return setting.scope[j].replace(/\./g, "_");
                    }
                }
            }
            var idx = scopes.lastIndexOf(".");
            var idx2 = scopes.lastIndexOf(" ");
            if (idx == idx2)
                break;
            scopes = scopes.slice(0, Math.max(idx, idx2));
        }
        //console.log("No scope found for " + scopes);
        return "default";
    }
    return this;
}

function SyntaxPattern(pattern)
{
    if (pattern.match)
    {
        this.match = new Regex(pattern.match);
    }
    if (pattern.begin)
    {
        this.begin = new Regex(pattern.begin);
    }
    if (pattern.end)
    {
        this.end = new Regex(pattern.end);
    }
    this.captures = pattern.captures;
    this.beginCaptures = pattern.beginCaptures;
    this.endCaptures = pattern.endCaptures;
    if (pattern.patterns)
    {
        this.patterns = pattern.patterns;
        for (var i in pattern.patterns)
        {
            this.patterns[i] = new SyntaxPattern(pattern.patterns[i]);
        }
    }
    this.name = pattern.name;
    return this;
}

function Syntax(name)
{
    var tmLang = loadFile(name);
    var jsonString = PlistParser.parse(toXML(tmLang));

    var patterns = jsonString.patterns;
    for (var i in patterns)
    {
        var pattern = patterns[i];
        patterns[i] = new SyntaxPattern(pattern);
    }
    this.jsonData = jsonString;

    this.firstMatch = function(data, patterns)
    {
        // Find the pattern that is the earliest match
        // TODO: this could be optimized
        var match = null;
        var startIdx = -1;
        var pattern = null;
        for (var i in patterns)
        {
            var innerPattern = patterns[i];
            if (innerPattern.match)
            {
                var innermatch = innerPattern.match.exec(data);
                if (innermatch)
                {
                    var idx = innermatch.index;
                    if (startIdx < 0 || startIdx > idx)
                    {
                        startIdx = idx;
                        match = innermatch;
                        pattern = innerPattern;
                        // console.log("" + pattern + ", match: " + match + ", idx: " + idx);
                    }
                }
            }
            else if (innerPattern.begin)
            {
                var innermatch = innerPattern.begin.exec(data);
                // TODO: remove duplicate..
                if (innermatch)
                {
                    var idx = innermatch.index;
                    if (startIdx < 0 || startIdx > idx)
                    {
                        startIdx = idx;
                        match = innermatch;
                        pattern = innerPattern;
                        // console.log("" + pattern + ", match: " + match + ", idx: " + idx);

                    }
                }
            }
        }
        return pattern;
    }

    this.innerApplyPattern = function(data, scope, match, captures)
    {
        var ret = "";
        if (match[0].indexOf("\"") != -1)
        {
            console.log(captures);
        }
        if (captures)
        {
            var lastIdx = 0;
            if (captures[0])
            {
                ret += "<span class=\"" + theme.getCssClassesForScopes(scope + " " + captures[0].name) + "\">";
            }

            for (var i = 1; i < match.length; i++)
            {
                if (!match[i])
                {
                    continue;
                }
                if (!match[0].slice(lastIdx).startsWith(match[i]))
                {
                    ret += match[0].slice(lastIdx, match[0].indexOf(match[i], lastIdx));
                }

                var capture = captures[i];
                var span = htmlify(match[i]);
                if (capture)
                {
                    span = "<span class=\"" + theme.getCssClassesForScopes(scope + " " + capture.name) + "\">" + span + "</span>";
                }

                ret += span;
                lastIdx = match[0].indexOf(match[i], lastIdx) + match[i].length;
            }
            if (lastIdx != match[0].length)
            {
                ret += match[0].slice(lastIdx);
            }
            if (captures[0])
            {
                ret += "</span>";
            }
        }
        else
        {
            ret += htmlify(match[0]);
        }
        fullline = match[0];
        start = match.index;

        var idx = start + fullline.length;
        data = data.slice(idx);
        return {"ret": ret, "data": data};
    }

    this.applyPattern = function(data, scope, pattern, theme)
    {
        var ret = "";
        var match = null;
        var start = 0;


        scope += " " + pattern.name;
        if (pattern.match)
        {
            match = pattern.match.exec(data);
        }
        else
        {
            match = pattern.begin.exec(data);
        }

        ret += htmlify(data.slice(0, match.index));
        ret += "<span class=\"" + theme.getCssClassesForScopes(scope) + "\">";
        var fullline = "";


        if (pattern.match)
        {
            var appl = this.innerApplyPattern(data, scope, match, pattern.captures);
            data = appl.data;
            ret += appl.ret;
        }
        else
        {
            match = pattern.begin.exec(data);
            var appl = this.innerApplyPattern(data, scope, match, pattern.beginCaptures)
            data = appl.data;
            ret += appl.ret;

            start = 0;

            var idx = start;
            var end = data.length;
            if (pattern.end)
            {
                while (data.length)
                {
                    var slice = data.slice(idx);
                    var match2 = pattern.end.exec(slice);
                    if (match2)
                    {
                        end = match2.index + idx + match2[0].length;
                    }

                    if (pattern.patterns)
                    {
                        var pattern2 = this.firstMatch(slice, pattern.patterns);
                        var match3 = null;

                        if (pattern2)
                        {
                            match3 = pattern2.match.exec(slice);
                        }
                        if (match3 && match3.index < match2.index)
                        {
                            var applied = this.applyPattern(slice, scope, pattern2, theme);
                            ret += applied.ret;
                            start = end = idx = 0;
                            data = applied.data;

                            continue;
                        }
                    }
                    ret += htmlify(data.slice(0, match2.index));
                    var appl = this.innerApplyPattern(slice, scope, match2, pattern.endCaptures)
                    data = appl.data;
                    ret += appl.ret;
                    start = end = idx = 0;

                    break;
                }
            }
            if (start != end)
            {
                var span = data.slice(start, end);
                ret += htmlify(span);
                fullline = span;
            }
        }
        ret += "</span>"

        var idx = start + fullline.length;
        data = data.slice(idx);
        return {"ret": ret, "data": data};
    }
    this.transform = function(data, theme)
    {
        var ret = "";
        ret += "<pre class=\"" + theme.getCssClassesForScopes(this.jsonData.scopeName) + "\">";

        var max = 10000;
        while (data.length > 0 && --max > 0)
        {
            var scope = this.jsonData.scopeName;
            var pattern = this.firstMatch(data, this.jsonData.patterns);

            if (!pattern)
            {
                // No more matches found
                break;
            }
            else
            {
                var applied = this.applyPattern(data, scope, pattern, theme);
                ret += applied.ret;
                data = applied.data;
            }
        }
        ret += "</pre>";
        return ret;
    }
    return this;
}

var theme = new Theme("/Packages/Color Scheme - Default/Monokai.tmTheme")
var syntax = new Syntax("Packages/JavaScript/JavaScript.tmLanguage");
var data = loadFile("lime.js");
var tdata = syntax.transform(data, theme);
document.write(tdata);
document.getElementsByTagName('body').innerHTML = tdata;

// console.log(syntax.transform("// test\nbice", theme));
// console.log(syntax.transform("// test\n", theme));
// console.log(syntax.transform("// test", theme));
