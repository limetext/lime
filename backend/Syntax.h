#ifndef __INCLUDED_LIME_BACKEND_SYNTAX_H
#define __INCLUDED_LIME_BACKEND_SYNTAX_H


#include <boost/regex.hpp>
#include <boost/any.hpp>
#include <vector>
#include <map>
#include <vector>

namespace lime
{
namespace backend
{

class Syntax;
class SyntaxPattern
{
public:
    typedef std::map<std::string, boost::any> Dict;
    typedef std::vector<boost::any> Array;

    SyntaxPattern(const Dict &dict, SyntaxPattern* root);
    ~SyntaxPattern();
private:
    boost::regex* mMatch;
    boost::regex* mBegin;
    boost::regex* mEnd;

    std::string mName;

    std::vector<boost::shared_ptr<SyntaxPattern> > mPatterns;
};

class Syntax
{
public:
    Syntax(const char* filename);
    ~Syntax();
private:
    SyntaxPattern* mRoot;
};


}

}


#endif
