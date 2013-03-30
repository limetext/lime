package textmate

import (
	"code.google.com/p/log4go"
	"encoding/json"
	"fmt"
	"github.com/moovweb/rubex"
	"github.com/quarnster/parser"
	"strconv"
	"strings"
)

const maxiter = 10000

type (
	Regex struct {
		re *rubex.Regexp
	}

	Language struct {
		UnpatchedLanguage
	}

	LanguageProvider interface {
		GetLanguage(id string) (*Language, error)
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
		Patterns       []*Pattern
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
		Language *Language
		root     parser.Node
	}
)

var (
	Provider LanguageProvider
	failed   = make(map[string]bool)
)

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

func (r *Regex) String() string {
	return r.re.String()
}

func (r *RootPattern) String() (ret string) {
	for i := range r.Patterns {
		ret += fmt.Sprintf("\t%s\n", r.Patterns[i])
	}
	return
}

func (s *Language) String() string {
	return fmt.Sprintf("%s\n%s", s.ScopeName, s.RootPattern)
}

func (p *Pattern) setOwner(l *Language) {
	p.owner = l
	for i := range p.Patterns {
		p.Patterns[i].setOwner(l)
	}
}

func (l *Language) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &l.UnpatchedLanguage); err != nil {
		return err
	}
	l.RootPattern.setOwner(l)
	for k := range l.Repository {
		l.Repository[k].setOwner(l)
	}
	return nil
}

func (r *RootPattern) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &r.Patterns)
}

func (r *Regex) UnmarshalJSON(data []byte) error {
	str := string(data[1 : len(data)-1])
	str = strings.Replace(str, "\\\\", "\\", -1)
	if re, err := rubex.Compile(str); err != nil {
		log4go.Warn("Couldn't compile language pattern %s: %s", str, err)
	} else {
		r.re = re
	}
	return nil
}

func (m MatchObject) fix(add int) {
	for i := range m {
		m[i] += add
	}
}

func (p *Pattern) FirstMatch(data string, pos int) (pat *Pattern, ret MatchObject) {
	startIdx := -1
	for i := 0; i < len(p.cachedPatterns); {
		ip, im := p.cachedPatterns[i].Cache(data, pos)
		if im != nil {
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
		if p.cachedMatch[0] >= pos {
			p.hits++
			return p.cachedPat, p.cachedMatch
		}
	} else {
		p.cachedPatterns = nil
	}
	if p.cachedPatterns == nil {
		p.cachedPatterns = make([]*Pattern, len(p.Patterns))
		copy(p.cachedPatterns, p.Patterns)
	}
	p.misses++

	if p.Match.re != nil {
		pat, ret = p, p.Match.re.FindStringSubmatchIndex(data[pos:])
		ret.fix(pos)
	} else if p.Begin.re != nil {
		pat, ret = p, p.Begin.re.FindStringSubmatchIndex(data[pos:])
		ret.fix(pos)
	} else if p.Include != "" {
		if z := p.Include[0]; z == '#' {
			key := p.Include[1:]
			if p2, ok := p.owner.Repository[key]; ok {
				pat, ret = p2.Cache(data, pos)
			} else {
				log4go.Error("Not found in repository: %s", p.Include)
			}
		} else if z == '$' {
			// TODO(q): Implement tmLanguage $ include directives
			log4go.Warn("Unhandled include directive: %s", p.Include)
		} else if Provider == nil {
			log4go.Warn("Include directive %s couldn't be handled as there's no Provider set", p.Include)
		} else if l, err := Provider.GetLanguage(p.Include); err != nil {
			if !failed[p.Include] {
				log4go.Error("Include directive %s failed: %s", p.Include, err)
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
			child := &parser.Node{Name: v.Name, Range: ranges[i], P: d}
			parents[i] = child
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
	fmt.Println("skipping: ", data[pos:mo[0]], "consuming:", data[mo[0]:mo[1]])
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
				endmatch := MatchObject(p.End.re.FindStringSubmatchIndex(data[i:]))
				endmatch.fix(i)
				//fmt.Println(i, end, endmatch, p.End.re)
				if endmatch != nil {
					end = endmatch[1]
				} else {
					if !found {
						// oops.. no end found at all, set it to the next line
						if e2 := strings.IndexRune(data[i:], '\n'); e2 != -1 {
							end = i + e2
						}
						//	fmt.Println("break 1")
						break
					} else {
						end = i
						//fmt.Println("break 2")
						break
					}
				}
				if (endmatch == nil || (endmatch != nil && endmatch[0] != i+1)) && len(p.cachedPatterns) > 0 {
					// Might be more recursive patterns to apply BEFORE the end is reached
					pattern2, match2 := p.FirstMatch(data, i)
					if match2 != nil &&
						((endmatch == nil && match2[0] < end) ||
							(endmatch != nil && match2[0] < endmatch[0])) {

						found = true
						r := pattern2.CreateNode(data, i, d, match2)
						ret.Append(r)
						i = r.Range.End
						continue
					}
				}
				if endmatch != nil {
					p.CreateCaptureNodes(data, i, d, endmatch, &ret, p.EndCaptures)
				}
				//fmt.Println("break 3")
				break
			}
			ret.Range.End = end
		}
	}
	ret.UpdateRange()
	return &ret
}

type dp struct {
	data string
}

func (d *dp) Data(a, b int) string {
	return d.data[a:b]
}

func (lp *LanguageParser) Parse(data string) bool {
	d := &dp{data}
	lp.root = parser.Node{P: d}
	iter := maxiter
	for i := 0; i < len(data) && iter > 0; iter-- {
		fmt.Println(i, len(data))
		pat, ret := lp.Language.RootPattern.Cache(data, i)
		if ret == nil {
			break
		} else {
			n := pat.CreateNode(data, i, d, ret)
			lp.root.Append(n)

			i = n.Range.End
		}
	}
	lp.root.UpdateRange()
	fmt.Println(lp.root)
	return true
}
