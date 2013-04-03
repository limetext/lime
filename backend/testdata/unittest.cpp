#include <vector>

class A
{
public:
    int test;
};

A a;
a.test;

std::vector<A> v;

v.front().test;


typedef std::vector<A> AV;

AV av;
av.front().test;

class B
{
public:
    AV variable;

    std::vector<A> variable2;
};

B b;
b.variable.front().test;
b.variable2.front().test;


class B2
{
public:
    std::vector<AV> variable;
    std::vector<std::vector<A> > variable2;
};

B2 b2;
b2.variable.front().back().test;
b2.variable2.front().back().test;


class B3
{
public:
    std::vector<std::vector<AV> > variable;
    std::vector<std::vector<std::vector<A> > > variable2;
};

B3 b3;
b3.variable.front().back().front().test;
b3.variable2.front().back().front().test;


template <typename T>
class TempA
{
public:
    T& GetT() { return mT;}
    T mT;
};

template <typename T>
class TempB
{
public:
    T& GetTB() { return mT; }
    T mTB;
};

TempA<A> ta;
ta.GetT().test;

TempA<TempB<TempA<TempB<A> > > > ta2;
ta2.GetT().GetTB().GetT().GetTB().test;
ta2.mT.mTB.mT.mTB.test;

typedef TempA<TempB<TempA<TempB<A> > > > Tababa;

Tababa tababa;
tababa.GetT().GetTB().GetT().GetTB().test;
tababa.mT.mTB.mT.mTB.test;

class Tababa_class
{
public:
    Tababa GetTababa();
    Tababa mTababa;

    TempA<TempB<TempA<TempB<A> > > > GetTababa2();
    TempA<TempB<TempA<TempB<A> > > > mTababa2;
};

Tababa_class tababa2;
tababa2.GetTababa().GetT().GetTB().GetT().GetTB().test;
tababa2.mTababa.GetT().GetTB().GetT().GetTB().test;


#include <string>
#include <boost/shared_ptr.hpp>

std::string str;
str.c_str();

boost::shared_ptr<std::vector<std::string> > strtest;
strtest.reset();
strtest->back()->c_str();

int test(AV& v)
{
    v.back().test;
}


