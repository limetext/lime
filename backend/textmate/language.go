package textmate

import (
	"code.google.com/p/log4go"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/quarnster/parser"
	"io/ioutil"
	"lime/3rdparty/libs/rubex"
	"lime/backend/loaders"
	"lime/backend/primitives"
	"strconv"
	"strings"
	"sync"
)

const maxiter = 10000

type (
	Regex struct {
		re        *rubex.Regexp
		lastIndex int
		lastFound int
	}

	Language struct {
		UnpatchedLanguage
	}

	LanguageProvider struct {
		sync.Mutex
		scope map[string]string
	}

	UnpatchedLanguage struct {
		FileTypes      []string
		FirstLineMatch string
		RootPattern    RootPattern `json:"patterns"`
		Repository     map[string]*Pattern
		ScopeName      string
	}

	Named struct {
		Name string
	}

	Captures map[string]Named

	MatchObject []int

	Pattern struct {
		Named
		Include        string
		Match          Regex
		Captures       Captures
		Begin          Regex
		BeginCaptures  Captures
		End            Regex
		EndCaptures    Captures
		Patterns       []Pattern
		owner          *Language // needed for include directives
		cachedData     string
		cachedPat      *Pattern
		cachedPatterns []*Pattern
		cachedMatch    MatchObject
		hits           int
		misses         int
	}
	RootPattern struct {
		Pattern
	}

	LanguageParser struct {
		l    *Language
		data []rune
	}
)

var (
	Provider LanguageProvider
	failed   = make(map[string]bool)
)

func init() {
	Provider.scope = make(map[string]string)
}

func (t *LanguageProvider) GetLanguage(id string) (*Language, error) {
	if l, err := t.LanguageFromScope(id); err != nil {
		return t.LanguageFromFile(id)
	} else {
		return l, err
	}
}

func (t *LanguageProvider) LanguageFromScope(id string) (*Language, error) {
	t.Lock()
	s, ok := t.scope[id]
	t.Unlock()
	if !ok {
		return nil, errors.New("Can't handle id " + id)
	} else {
		return t.LanguageFromFile(s)
	}
}

func (t *LanguageProvider) LanguageFromFile(fn string) (*Language, error) {
	if d, err := ioutil.ReadFile(fn); err != nil {
		return nil, fmt.Errorf("Couldn't load file %s: %s", fn, err)
	} else {
		var l Language
		if err := loaders.LoadPlist(d, &l); err != nil {
			return nil, err
		} else {
			t.Lock()
			defer t.Unlock()
			t.scope[l.ScopeName] = fn
			return &l, nil
		}
	}
}

func (p Pattern) String() (ret string) {
	ret = fmt.Sprintf(`---------------------------------------
Name:    %s
Match:   %s
Begin:   %s
End:     %s
Include: %s
`, p.Name, p.Match, p.Begin, p.End, p.Include)
	ret += fmt.Sprintf("<Sub-Patterns>\n")
	for i := range p.Patterns {
		inner := fmt.Sprintf("%s", p.Patterns[i])
		ret += fmt.Sprintf("\t%s\n", strings.Replace(strings.Replace(inner, "\t", "\t\t", -1), "\n", "\n\t", -1))
	}
	ret += fmt.Sprintf("</Sub-Patterns>\n---------------------------------------")
	return
}

func (r Regex) String() string {
	if r.re == nil {
		return "nil"
	}
	return fmt.Sprintf("%s   // %d, %d", r.re.String(), r.lastIndex, r.lastFound)
}

func (r *RootPattern) String() (ret string) {
	for i := range r.Patterns {
		ret += fmt.Sprintf("\t%s\n", r.Patterns[i])
	}
	return
}

func (s *Language) String() string {
	return fmt.Sprintf("%s\n%s\n", s.ScopeName, s.RootPattern, s.Repository)
}

func (p *Pattern) tweak(l *Language) {
	p.owner = l
	p.Name = strings.TrimSpace(p.Name)
	for i := range p.Patterns {
		p.Patterns[i].tweak(l)
	}
}

func (l *Language) tweak() {
	l.RootPattern.tweak(l)
	for k := range l.Repository {
		p := l.Repository[k]
		p.tweak(l)
		l.Repository[k] = p
	}
}

func (l *Language) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &l.UnpatchedLanguage); err != nil {
		return err
	}
	l.tweak()
	return nil
}

func (r *RootPattern) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Patterns)
}

func (r *Regex) UnmarshalJSON(data []byte) error {
	str := string(data[1 : len(data)-1])
	str = strings.Replace(str, "\\\\", "\\", -1)
	str = strings.Replace(str, "\\n", "\n", -1)
	str = strings.Replace(str, "\\t", "\t", -1)
	if re, err := rubex.Compile(str); err != nil {
		log4go.Warn("Couldn't compile language pattern %s: %s", str, err)
	} else {
		r.re = re
	}
	return nil
}

func (m MatchObject) fix(add int) {
	for i := range m {
		if m[i] != -1 {
			m[i] += add
		}
	}
}

func (r *Regex) Find(data string, pos int) MatchObject {
	if r.lastIndex > pos {
		r.lastFound = 0
	}
	r.lastIndex = pos
	for r.lastFound < len(data) {
		ret := r.re.FindStringSubmatchIndex(data[r.lastFound:])
		if ret == nil {
			break
		} else if (ret[0] + r.lastFound) < pos {
			if ret[0] == 0 {
				r.lastFound++
			} else {
				r.lastFound += ret[0]
			}
			continue
		}
		mo := MatchObject(ret)
		mo.fix(r.lastFound)
		return mo
	}
	return nil
}

func (p *Pattern) FirstMatch(data string, pos int) (pat *Pattern, ret MatchObject) {
	startIdx := -1
	for i := 0; i < len(p.cachedPatterns); {
		ip, im := p.cachedPatterns[i].Cache(data, pos)
		if im != nil /* && im[0] != im[1]*/ {
			if startIdx < 0 || startIdx > im[0] {
				startIdx, pat, ret = im[0], ip, im
				// This match is right at the start, we're not going to find a better pattern than this,
				// so stop the search
				if im[0] == pos {
					break
				}
			}
			i++
		} else {
			// If it wasn't found now, it'll never be found, so the pattern can be popped from the cache
			copy(p.cachedPatterns[i:], p.cachedPatterns[i+1:])
			p.cachedPatterns = p.cachedPatterns[:len(p.cachedPatterns)-1]
		}
	}
	return
}

func (p *Pattern) Cache(data string, pos int) (pat *Pattern, ret MatchObject) {
	if p.cachedData == data {
		if p.cachedMatch == nil {
			return nil, nil
		}
		if p.cachedMatch[0] >= pos && p.cachedPat.cachedMatch != nil {
			p.hits++
			return p.cachedPat, p.cachedMatch
		}
	} else {
		p.cachedPatterns = nil
	}
	if p.cachedPatterns == nil {
		p.cachedPatterns = make([]*Pattern, len(p.Patterns))
		for i := range p.cachedPatterns {
			p.cachedPatterns[i] = &p.Patterns[i]
		}
	}
	p.misses++

	if p.Match.re != nil {
		pat, ret = p, p.Match.Find(data, pos)
	} else if p.Begin.re != nil {
		pat, ret = p, p.Begin.Find(data, pos)
	} else if p.Include != "" {
		if z := p.Include[0]; z == '#' {
			key := p.Include[1:]
			if p2, ok := p.owner.Repository[key]; ok {
				pat, ret = p2.Cache(data, pos)
			} else {
				log4go.Warn("Not found in repository: %s", p.Include)
			}
		} else if z == '$' {
			// TODO(q): Implement tmLanguage $ include directives
			log4go.Warn("Unhandled include directive: %s", p.Include)
		} else if l, err := Provider.GetLanguage(p.Include); err != nil {
			if !failed[p.Include] {
				log4go.Warn("Include directive %s failed: %s", p.Include, err)
			}
			failed[p.Include] = true
		} else {
			return l.RootPattern.Cache(data, pos)
		}
	} else {
		pat, ret = p.FirstMatch(data, pos)
	}
	p.cachedData = data
	p.cachedMatch = ret
	p.cachedPat = pat

	return
}

func (p *Pattern) CreateCaptureNodes(data string, pos int, d parser.DataSource, mo MatchObject, parent *parser.Node, cap Captures) {
	ranges := make([]parser.Range, len(mo)/2)
	parentIndex := make([]int, len(ranges))
	parents := make([]*parser.Node, len(parentIndex))
	for i := range ranges {
		ranges[i] = parser.Range{mo[i*2+0], mo[i*2+1]}
		if i < 2 {
			parents[i] = parent
			continue
		}
		r := ranges[i]
		for j := i - 1; j >= 0; j-- {
			if ranges[j].Contains(r) {
				parentIndex[i] = j
				break
			}
		}
	}

	for k, v := range cap {
		i64, err := strconv.ParseInt(k, 10, 32)
		if i := int(i64); err == nil && i < len(parents) {
			if ranges[i].Start == -1 {
				continue
			}
			child := &parser.Node{Name: v.Name, Range: ranges[i], P: d}
			parents[i] = child
			if i == 0 {
				parent.Append(child)
				continue
			}
			var p *parser.Node
			for p == nil {
				i = parentIndex[i]
				p = parents[i]
			}
			p.Append(child)
		}
	}
}

func (p *Pattern) CreateNode(data string, pos int, d parser.DataSource, mo MatchObject) *parser.Node {
	ret := parser.Node{Name: p.Name, Range: parser.Range{mo[0], mo[1]}, P: d}
	if p.Match.re != nil {
		p.CreateCaptureNodes(data, pos, d, mo, &ret, p.Captures)
	} else if p.Begin.re != nil {
		if len(p.BeginCaptures) > 0 {
			p.CreateCaptureNodes(data, pos, d, mo, &ret, p.BeginCaptures)
		} else {
			p.CreateCaptureNodes(data, pos, d, mo, &ret, p.Captures)
		}

		if p.End.re != nil {
			var (
				found  = false
				i, end int
			)
			for i, end = ret.Range.End, len(data); i < len(data); {
				endmatch := p.End.Find(data, i)
				if endmatch != nil {
					end = endmatch[1]
				} else {
					if !found {
						// oops.. no end found at all, set it to the next line
						if e2 := strings.IndexRune(data[i:], '\n'); e2 != -1 {
							end = i + e2
						} else {
							end = len(data)
						}
						break
					} else {
						end = i
						break
					}
				}
				if /*(endmatch == nil || (endmatch != nil && endmatch[0] != i)) && */ len(p.cachedPatterns) > 0 {
					// Might be more recursive patterns to apply BEFORE the end is reached
					pattern2, match2 := p.FirstMatch(data, i)
					if match2 != nil && ((endmatch == nil && match2[0] < end) || (endmatch != nil && (match2[0] < endmatch[0] || match2[0] == endmatch[0] && ret.Range.Start == ret.Range.End))) {
						found = true
						r := pattern2.CreateNode(data, i, d, match2)
						ret.Append(r)
						i = r.Range.End
						continue
					}
				}
				if endmatch != nil {
					if len(p.EndCaptures) > 0 {
						p.CreateCaptureNodes(data, i, d, endmatch, &ret, p.EndCaptures)
					} else {
						p.CreateCaptureNodes(data, i, d, endmatch, &ret, p.Captures)
					}
				}
				break
			}
			ret.Range.End = end
		}
	}
	ret.UpdateRange()
	return &ret
}

func (d *LanguageParser) Data(a, b int) string {
	a = primitives.Clamp(0, len(d.data), a)
	b = primitives.Clamp(0, len(d.data), b)
	return string(d.data[a:b])
}

func (lp *LanguageParser) patch(lut []int, node *parser.Node) {
	node.Range.Start = lut[node.Range.Start]
	node.Range.End = lut[node.Range.End]
	for _, child := range node.Children {
		lp.patch(lut, child)
	}
}

func NewLanguageParser(scope string, data string) (*LanguageParser, error) {
	if l, err := Provider.GetLanguage(scope); err != nil {
		return nil, err
	} else {
		return &LanguageParser{l, []rune(data)}, nil
	}
}

func (lp *LanguageParser) Parse() (*parser.Node, error) {
	sdata := string(lp.data)
	rn := parser.Node{P: lp, Name: lp.l.ScopeName}
	defer func() {
		if r := recover(); r != nil {
			log4go.Error("Panic during parse: %v\n", r)
			log4go.Debug("%v", rn)
		}
	}()
	iter := maxiter
	for i := 0; i < len(sdata) && iter > 0; iter-- {
		pat, ret := lp.l.RootPattern.Cache(sdata, i)
		nl := strings.IndexAny(sdata[i:], "\n\r")
		if nl != -1 {
			nl += i
		}
		if ret == nil {
			break
		} else if nl > 0 && nl <= ret[0] {
			i = nl
			for i < len(sdata) && (sdata[i] == '\n' || sdata[i] == '\r') {
				i++
			}
		} else {
			n := pat.CreateNode(sdata, i, lp, ret)
			rn.Append(n)

			i = n.Range.End
		}
	}
	rn.UpdateRange()
	if len(sdata) != 0 {
		lut := make([]int, len(sdata)+1)
		j := 0
		for i := range sdata {
			lut[i] = j
			j++
		}
		lut[len(sdata)] = len(lp.data)
		lp.patch(lut, &rn)
	}
	if iter == 0 {
		panic("reached maximum number of iterations")
	}
	return &rn, nil
}
