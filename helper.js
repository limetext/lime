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

// This XRegExp addon code is from: https://gist.github.com/2387872
// MIT Licensed
(function (XRegExp) {

    function prepareLb(lb) {
        // Allow mode modifier before lookbehind
        var parts = /^((?:\(\?[\w$]+\))?)\(\?<([=!])([\s\S]*)\)$/.exec(lb);
        return {
            // $(?!\s) allows use of (?m) in lookbehind
            lb: XRegExp(parts ? parts[1] + "(?:" + parts[3] + ")$(?!\\s)" : lb),
            // Positive or negative lookbehind. Use positive if no lookbehind group
            type: parts ? parts[2] === "=" : !parts
        };
    }

    XRegExp.execLb = function (str, lb, regex) {
        var pos = 0, match, leftContext;
        lb = prepareLb(lb);
        while (match = XRegExp.exec(str, regex, pos)) {
            leftContext = str.slice(0, match.index);
            if (lb.type === lb.lb.test(leftContext)) {
                return match;
            }
            pos = match.index + 1;
        }
        return null;
    };

    XRegExp.testLb = function (str, lb, regex) {
        return !!XRegExp.execLb(str, lb, regex);
    };

    XRegExp.searchLb = function (str, lb, regex) {
        var match = XRegExp.execLb(str, lb, regex);
        return match ? match.index : -1;
    };

    XRegExp.matchAllLb = function (str, lb, regex) {
        var matches = [], pos = 0, match, leftContext;
        lb = prepareLb(lb);
        while (match = XRegExp.exec(str, regex, pos)) {
            leftContext = str.slice(0, match.index);
            if (lb.type === lb.lb.test(leftContext)) {
                matches.push(match[0]);
                pos = match.index + (match[0].length || 1);
            } else {
                pos = match.index + 1;
            }
        }
        return matches;
    };

    XRegExp.replaceLb = function (str, lb, regex, replacement) {
        var output = "", pos = 0, lastEnd = 0, match, leftContext;
        lb = prepareLb(lb);
        while (match = XRegExp.exec(str, regex, pos)) {
            leftContext = str.slice(0, match.index);
            if (lb.type === lb.lb.test(leftContext)) {
                // Doesn't work correctly if lookahead in regex looks outside of the match
                output += str.slice(lastEnd, match.index) + XRegExp.replace(match[0], regex, replacement);
                lastEnd = match.index + match[0].length;
                if (!regex.global) {
                    break;
                }
                pos = match.index + (match[0].length || 1);
            } else {
                pos = match.index + 1;
            }
        }
        return output + str.slice(lastEnd);
    };

}(XRegExp));

if (typeof String.prototype.startsWith != 'function') {
    String.prototype.startsWith = function (str){
        return this.slice(0, str.length) == str;
    };
}
if (typeof String.prototype.endsWith != 'function') {
    String.prototype.endsWith = function (str){
        if (str.length > this.length)
        {
            return false;
        }
        return this.slice(this.length-str.length) == str;
    };
}

if (typeof String.prototype.trim != 'function') {
    String.prototype.trim = function()
    {
        return this.replace(/^\s\s*/, '').replace(/\s\s*$/, '');
    };
}



function Regex(pattern, flags)
{
    if (typeof(flags) === 'undefined')
        flags = "m";

    pattern = pattern.replace("\\h", "[a-fA-F0-9]");
    if (pattern.startsWith("(?<"))
    {
        var temp = XRegExp("(\\\([^\)]+\\\))(.*)").exec(pattern);
        this.lookback = temp[1];
        pattern = temp[2];
    }
    else
    {
        this.lookback = null;
    }

    if (pattern.indexOf("(?<") != -1)
    {
        console.log("Warning, lookback regex is used somewhere else than right at the start. Lookback pattern will be stripped");
        var pat = "(.*)(\\\(\\\?\\\<[^\)]+\\\))(.*)";
        var re = XRegExp(pat);
        var res = null;
        var pattern2 = pattern;
        while (res = re.exec(pattern2))
        {
            pattern2 = res[0].replace(res[2], "");
        }

        pattern = pattern2;
    }
    if (pattern.indexOf("\\G") != -1)
    {
        console.log("Warning, unsupported \\G used. \\G will be stripped");
        pattern = pattern.replace("\\G", "");
    }
    try
    {
        this.pattern = XRegExp(pattern, flags);
    }
    catch (e)
    {
        console.log("Warning, could't create regex pattern: " + e.message);
        this.pattern = null;
    }

    this.exec = function(str)
    {
        if (this.pattern)
        {
            if (this.lookback)
            {
                return XRegExp.execLb(str, this.lookback, this.pattern);
            }

            return this.pattern.exec(str);
        }
        return null;
    }
    return this;
}

// hexToRgb and rgbToHex from http://stackoverflow.com/questions/5623838/rgb-to-hex-and-hex-to-rgb
function hexToRgb(hex) {
    var result = /^#?(a-f\d){0,2}([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
    return result ? {
        r: parseInt(result[2], 16),
        g: parseInt(result[3], 16),
        b: parseInt(result[4], 16)
    } : null;
}

function rgbToHex(r, g, b)
{
    return "#" + ((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1);
}

function clamp(nb, min, max)
{
    return Math.max(min, Math.min(max, nb));
}

function htmlify(str)
{
    return str.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
}

function loadFile(name)
{
    var client = new XMLHttpRequest();
    client.open('GET', name, false);
    client.send();
    return client.responseText;
}

function toXML(text)
{
    var parser = new DOMParser();
    var doc = parser.parseFromString(text,'text/xml');
    return doc;
}
