#include "Syntax.h"

#include <map>
#include <string>
#include <boost/any.hpp>

#include <Plist.hpp>
#include <boost/foreach.hpp>
#include <boost/exception/diagnostic_information.hpp>

namespace lime
{
namespace backend
{

Syntax::Syntax(const char* filename)
{
    SyntaxPattern::Dict map;

    Plist::readPlist(filename, map);
    mRoot = new SyntaxPattern(map, NULL);
}

Syntax::~Syntax()
{
    delete mRoot;
}

boost::regex* CreateRegex(const SyntaxPattern::Dict &dict, const std::string& key)
{
    SyntaxPattern::Dict::const_iterator i = dict.find(key);
    boost::regex* ret = NULL;

    if (i != dict.end())
    {
        try
        {
            const std::string &str = boost::any_cast<std::string>(i->second);
            ret = new boost::regex(str, boost::regex::perl);
        }
        catch (boost::exception &e)
        {
            printf("%s\n", boost::diagnostic_information(e).c_str());
        }
    }

    return ret;
}

SyntaxPattern::SyntaxPattern(const SyntaxPattern::Dict &dict, SyntaxPattern* parent)
    : mMatch(NULL), mBegin(NULL), mEnd(NULL)
{
    mMatch = CreateRegex(dict, "match");
    mBegin = CreateRegex(dict, "begin");
    mEnd   = CreateRegex(dict, "end");

    SyntaxPattern::Dict::const_iterator i;
    if (!parent)
    {
        i = dict.find("scopeName");
    }
    else
    {
        i = dict.find("name");
    }

    if (i != dict.end())
    {
        mName = boost::any_cast<std::string>(i->second);
    }
    printf("name: %s\n", mName.c_str());

    i = dict.find("patterns");
    if (i != dict.end())
    {
        if (!parent)
            parent = this;

        BOOST_FOREACH(const boost::any& data, boost::any_cast<SyntaxPattern::Array>(i->second))
        {
            const SyntaxPattern::Dict &d = boost::any_cast<SyntaxPattern::Dict>(data);
            mPatterns.push_back(boost::shared_ptr<SyntaxPattern>(new SyntaxPattern(d, parent)));
        }
    }
}
SyntaxPattern::~SyntaxPattern()
{
    delete mMatch;
    delete mBegin;
    delete mEnd;
    mPatterns.clear();
}

}
}
