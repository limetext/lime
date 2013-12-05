// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// This file was generated, and shouldn't be manually changed
package plist

import (
	. "github.com/quarnster/parser"
	"github.com/quarnster/util/text"
)

type PLIST struct {
	ParserData  Reader
	IgnoreRange text.Region
	Root        Node
	LastError   int
}

func (p *PLIST) RootNode() *Node {
	return &p.Root
}

func (p *PLIST) SetData(data string) {
	p.ParserData = NewReader(data)
	p.Root = Node{Name: "PLIST", P: p}
	p.IgnoreRange = text.Region{}
	p.LastError = 0
}

func (p *PLIST) Parse(data string) bool {
	p.SetData(data)
	ret := p.realParse()
	p.Root.UpdateRange()
	return ret
}

func (p *PLIST) Data(start, end int) string {
	return p.ParserData.Substring(start, end)
}

func (p *PLIST) Error() Error {
	errstr := ""
	line, column := p.ParserData.LineCol(p.LastError)

	if p.LastError == p.ParserData.Len() {
		errstr = "Unexpected EOF"
	} else {
		p.ParserData.Seek(p.LastError)
		if r := p.ParserData.Read(); r == '\r' || r == '\n' {
			errstr = "Unexpected new line"
		} else {
			errstr = "Unexpected " + string(r)
		}
	}
	return NewError(line, column, errstr)
}

func (p *PLIST) realParse() bool {
	return p.PlistFile()
}
func (p *PLIST) PlistFile() bool {
	// PlistFile      <-    "<?xml" (!"?>" .)+ "?>" Spacing* "<!DOCTYPE" (!'>' .)+ '>' Spacing* Plist Spacing* EndOfFile?
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			accept = true
			s := p.ParserData.Pos()
			if p.ParserData.Read() != '<' || p.ParserData.Read() != '?' || p.ParserData.Read() != 'x' || p.ParserData.Read() != 'm' || p.ParserData.Read() != 'l' {
				p.ParserData.Seek(s)
				accept = false
			}
		}
		if accept {
			{
				save := p.ParserData.Pos()
				{
					save := p.ParserData.Pos()
					s := p.ParserData.Pos()
					{
						accept = true
						s := p.ParserData.Pos()
						if p.ParserData.Read() != '?' || p.ParserData.Read() != '>' {
							p.ParserData.Seek(s)
							accept = false
						}
					}
					p.ParserData.Seek(s)
					p.Root.Discard(s)
					accept = !accept
					if accept {
						if p.ParserData.Pos() >= p.ParserData.Len() {
							accept = false
						} else {
							p.ParserData.Read()
							accept = true
						}
						if accept {
						}
					}
					if !accept {
						if p.LastError < p.ParserData.Pos() {
							p.LastError = p.ParserData.Pos()
						}
						p.ParserData.Seek(save)
					}
				}
				if !accept {
					p.ParserData.Seek(save)
				} else {
					for accept {
						{
							save := p.ParserData.Pos()
							s := p.ParserData.Pos()
							{
								accept = true
								s := p.ParserData.Pos()
								if p.ParserData.Read() != '?' || p.ParserData.Read() != '>' {
									p.ParserData.Seek(s)
									accept = false
								}
							}
							p.ParserData.Seek(s)
							p.Root.Discard(s)
							accept = !accept
							if accept {
								if p.ParserData.Pos() >= p.ParserData.Len() {
									accept = false
								} else {
									p.ParserData.Read()
									accept = true
								}
								if accept {
								}
							}
							if !accept {
								if p.LastError < p.ParserData.Pos() {
									p.LastError = p.ParserData.Pos()
								}
								p.ParserData.Seek(save)
							}
						}
					}
					accept = true
				}
			}
			if accept {
				{
					accept = true
					s := p.ParserData.Pos()
					if p.ParserData.Read() != '?' || p.ParserData.Read() != '>' {
						p.ParserData.Seek(s)
						accept = false
					}
				}
				if accept {
					{
						accept = true
						for accept {
							accept = p.Spacing()
						}
						accept = true
					}
					if accept {
						{
							accept = true
							s := p.ParserData.Pos()
							if p.ParserData.Read() != '<' || p.ParserData.Read() != '!' || p.ParserData.Read() != 'D' || p.ParserData.Read() != 'O' || p.ParserData.Read() != 'C' || p.ParserData.Read() != 'T' || p.ParserData.Read() != 'Y' || p.ParserData.Read() != 'P' || p.ParserData.Read() != 'E' {
								p.ParserData.Seek(s)
								accept = false
							}
						}
						if accept {
							{
								save := p.ParserData.Pos()
								{
									save := p.ParserData.Pos()
									s := p.ParserData.Pos()
									if p.ParserData.Read() != '>' {
										p.ParserData.UnRead()
										accept = false
									} else {
										accept = true
									}
									p.ParserData.Seek(s)
									p.Root.Discard(s)
									accept = !accept
									if accept {
										if p.ParserData.Pos() >= p.ParserData.Len() {
											accept = false
										} else {
											p.ParserData.Read()
											accept = true
										}
										if accept {
										}
									}
									if !accept {
										if p.LastError < p.ParserData.Pos() {
											p.LastError = p.ParserData.Pos()
										}
										p.ParserData.Seek(save)
									}
								}
								if !accept {
									p.ParserData.Seek(save)
								} else {
									for accept {
										{
											save := p.ParserData.Pos()
											s := p.ParserData.Pos()
											if p.ParserData.Read() != '>' {
												p.ParserData.UnRead()
												accept = false
											} else {
												accept = true
											}
											p.ParserData.Seek(s)
											p.Root.Discard(s)
											accept = !accept
											if accept {
												if p.ParserData.Pos() >= p.ParserData.Len() {
													accept = false
												} else {
													p.ParserData.Read()
													accept = true
												}
												if accept {
												}
											}
											if !accept {
												if p.LastError < p.ParserData.Pos() {
													p.LastError = p.ParserData.Pos()
												}
												p.ParserData.Seek(save)
											}
										}
									}
									accept = true
								}
							}
							if accept {
								if p.ParserData.Read() != '>' {
									p.ParserData.UnRead()
									accept = false
								} else {
									accept = true
								}
								if accept {
									{
										accept = true
										for accept {
											accept = p.Spacing()
										}
										accept = true
									}
									if accept {
										accept = p.Plist()
										if accept {
											{
												accept = true
												for accept {
													accept = p.Spacing()
												}
												accept = true
											}
											if accept {
												accept = p.EndOfFile()
												accept = true
												if accept {
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
		if !accept {
			if p.LastError < p.ParserData.Pos() {
				p.LastError = p.ParserData.Pos()
			}
			p.ParserData.Seek(save)
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) Plist() bool {
	// Plist          <-    "<plist version=\"1.0\">" Values "</plist>"
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			accept = true
			s := p.ParserData.Pos()
			if p.ParserData.Read() != '<' || p.ParserData.Read() != 'p' || p.ParserData.Read() != 'l' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 's' || p.ParserData.Read() != 't' || p.ParserData.Read() != ' ' || p.ParserData.Read() != 'v' || p.ParserData.Read() != 'e' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 's' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 'o' || p.ParserData.Read() != 'n' || p.ParserData.Read() != '=' || p.ParserData.Read() != '"' || p.ParserData.Read() != '1' || p.ParserData.Read() != '.' || p.ParserData.Read() != '0' || p.ParserData.Read() != '"' || p.ParserData.Read() != '>' {
				p.ParserData.Seek(s)
				accept = false
			}
		}
		if accept {
			accept = p.Values()
			if accept {
				{
					accept = true
					s := p.ParserData.Pos()
					if p.ParserData.Read() != '<' || p.ParserData.Read() != '/' || p.ParserData.Read() != 'p' || p.ParserData.Read() != 'l' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 's' || p.ParserData.Read() != 't' || p.ParserData.Read() != '>' {
						p.ParserData.Seek(s)
						accept = false
					}
				}
				if accept {
				}
			}
		}
		if !accept {
			if p.LastError < p.ParserData.Pos() {
				p.LastError = p.ParserData.Pos()
			}
			p.ParserData.Seek(save)
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) Dictionary() bool {
	// Dictionary     <-    "<dict>" KeyValuePair+ "</dict>" / "<dict/>"
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			save := p.ParserData.Pos()
			{
				accept = true
				s := p.ParserData.Pos()
				if p.ParserData.Read() != '<' || p.ParserData.Read() != 'd' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 'c' || p.ParserData.Read() != 't' || p.ParserData.Read() != '>' {
					p.ParserData.Seek(s)
					accept = false
				}
			}
			if accept {
				{
					save := p.ParserData.Pos()
					accept = p.KeyValuePair()
					if !accept {
						p.ParserData.Seek(save)
					} else {
						for accept {
							accept = p.KeyValuePair()
						}
						accept = true
					}
				}
				if accept {
					{
						accept = true
						s := p.ParserData.Pos()
						if p.ParserData.Read() != '<' || p.ParserData.Read() != '/' || p.ParserData.Read() != 'd' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 'c' || p.ParserData.Read() != 't' || p.ParserData.Read() != '>' {
							p.ParserData.Seek(s)
							accept = false
						}
					}
					if accept {
					}
				}
			}
			if !accept {
				if p.LastError < p.ParserData.Pos() {
					p.LastError = p.ParserData.Pos()
				}
				p.ParserData.Seek(save)
			}
		}
		if !accept {
			{
				accept = true
				s := p.ParserData.Pos()
				if p.ParserData.Read() != '<' || p.ParserData.Read() != 'd' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 'c' || p.ParserData.Read() != 't' || p.ParserData.Read() != '/' || p.ParserData.Read() != '>' {
					p.ParserData.Seek(s)
					accept = false
				}
			}
			if !accept {
			}
		}
		if !accept {
			p.ParserData.Seek(save)
		}
	}
	end := p.ParserData.Pos()
	if accept {
		node := p.Root.Cleanup(start, end)
		node.Name = "Dictionary"
		node.P = p
		node.Range = node.Range.Clip(p.IgnoreRange)
		p.Root.Append(node)
	} else {
		p.Root.Discard(start)
	}
	if p.IgnoreRange.A >= end || p.IgnoreRange.B <= start {
		p.IgnoreRange = text.Region{}
	}
	return accept
}

func (p *PLIST) KeyValuePair() bool {
	// KeyValuePair   <-    Spacing* KeyTag Spacing* Value Spacing*
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			accept = true
			for accept {
				accept = p.Spacing()
			}
			accept = true
		}
		if accept {
			accept = p.KeyTag()
			if accept {
				{
					accept = true
					for accept {
						accept = p.Spacing()
					}
					accept = true
				}
				if accept {
					accept = p.Value()
					if accept {
						{
							accept = true
							for accept {
								accept = p.Spacing()
							}
							accept = true
						}
						if accept {
						}
					}
				}
			}
		}
		if !accept {
			if p.LastError < p.ParserData.Pos() {
				p.LastError = p.ParserData.Pos()
			}
			p.ParserData.Seek(save)
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) KeyTag() bool {
	// KeyTag         <-    "<key>" Key "</key>"
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			accept = true
			s := p.ParserData.Pos()
			if p.ParserData.Read() != '<' || p.ParserData.Read() != 'k' || p.ParserData.Read() != 'e' || p.ParserData.Read() != 'y' || p.ParserData.Read() != '>' {
				p.ParserData.Seek(s)
				accept = false
			}
		}
		if accept {
			accept = p.Key()
			if accept {
				{
					accept = true
					s := p.ParserData.Pos()
					if p.ParserData.Read() != '<' || p.ParserData.Read() != '/' || p.ParserData.Read() != 'k' || p.ParserData.Read() != 'e' || p.ParserData.Read() != 'y' || p.ParserData.Read() != '>' {
						p.ParserData.Seek(s)
						accept = false
					}
				}
				if accept {
				}
			}
		}
		if !accept {
			if p.LastError < p.ParserData.Pos() {
				p.LastError = p.ParserData.Pos()
			}
			p.ParserData.Seek(save)
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) Key() bool {
	// Key            <-    (!'<' .)*
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		accept = true
		for accept {
			{
				save := p.ParserData.Pos()
				s := p.ParserData.Pos()
				if p.ParserData.Read() != '<' {
					p.ParserData.UnRead()
					accept = false
				} else {
					accept = true
				}
				p.ParserData.Seek(s)
				p.Root.Discard(s)
				accept = !accept
				if accept {
					if p.ParserData.Pos() >= p.ParserData.Len() {
						accept = false
					} else {
						p.ParserData.Read()
						accept = true
					}
					if accept {
					}
				}
				if !accept {
					if p.LastError < p.ParserData.Pos() {
						p.LastError = p.ParserData.Pos()
					}
					p.ParserData.Seek(save)
				}
			}
		}
		accept = true
	}
	end := p.ParserData.Pos()
	if accept {
		node := p.Root.Cleanup(start, end)
		node.Name = "Key"
		node.P = p
		node.Range = node.Range.Clip(p.IgnoreRange)
		p.Root.Append(node)
	} else {
		p.Root.Discard(start)
	}
	if p.IgnoreRange.A >= end || p.IgnoreRange.B <= start {
		p.IgnoreRange = text.Region{}
	}
	return accept
}

func (p *PLIST) StringTag() bool {
	// StringTag      <-    "<string>" String "</string>"
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			accept = true
			s := p.ParserData.Pos()
			if p.ParserData.Read() != '<' || p.ParserData.Read() != 's' || p.ParserData.Read() != 't' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 'n' || p.ParserData.Read() != 'g' || p.ParserData.Read() != '>' {
				p.ParserData.Seek(s)
				accept = false
			}
		}
		if accept {
			accept = p.String()
			if accept {
				{
					accept = true
					s := p.ParserData.Pos()
					if p.ParserData.Read() != '<' || p.ParserData.Read() != '/' || p.ParserData.Read() != 's' || p.ParserData.Read() != 't' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'i' || p.ParserData.Read() != 'n' || p.ParserData.Read() != 'g' || p.ParserData.Read() != '>' {
						p.ParserData.Seek(s)
						accept = false
					}
				}
				if accept {
				}
			}
		}
		if !accept {
			if p.LastError < p.ParserData.Pos() {
				p.LastError = p.ParserData.Pos()
			}
			p.ParserData.Seek(save)
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) String() bool {
	// String         <-    (!'<' .)*
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		accept = true
		for accept {
			{
				save := p.ParserData.Pos()
				s := p.ParserData.Pos()
				if p.ParserData.Read() != '<' {
					p.ParserData.UnRead()
					accept = false
				} else {
					accept = true
				}
				p.ParserData.Seek(s)
				p.Root.Discard(s)
				accept = !accept
				if accept {
					if p.ParserData.Pos() >= p.ParserData.Len() {
						accept = false
					} else {
						p.ParserData.Read()
						accept = true
					}
					if accept {
					}
				}
				if !accept {
					if p.LastError < p.ParserData.Pos() {
						p.LastError = p.ParserData.Pos()
					}
					p.ParserData.Seek(save)
				}
			}
		}
		accept = true
	}
	end := p.ParserData.Pos()
	if accept {
		node := p.Root.Cleanup(start, end)
		node.Name = "String"
		node.P = p
		node.Range = node.Range.Clip(p.IgnoreRange)
		p.Root.Append(node)
	} else {
		p.Root.Discard(start)
	}
	if p.IgnoreRange.A >= end || p.IgnoreRange.B <= start {
		p.IgnoreRange = text.Region{}
	}
	return accept
}

func (p *PLIST) Value() bool {
	// Value          <-    Array / StringTag / Dictionary
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		accept = p.Array()
		if !accept {
			accept = p.StringTag()
			if !accept {
				accept = p.Dictionary()
				if !accept {
				}
			}
		}
		if !accept {
			p.ParserData.Seek(save)
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) Values() bool {
	// Values         <-    (Spacing* Value Spacing*)*
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		accept = true
		for accept {
			{
				save := p.ParserData.Pos()
				{
					accept = true
					for accept {
						accept = p.Spacing()
					}
					accept = true
				}
				if accept {
					accept = p.Value()
					if accept {
						{
							accept = true
							for accept {
								accept = p.Spacing()
							}
							accept = true
						}
						if accept {
						}
					}
				}
				if !accept {
					if p.LastError < p.ParserData.Pos() {
						p.LastError = p.ParserData.Pos()
					}
					p.ParserData.Seek(save)
				}
			}
		}
		accept = true
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) Array() bool {
	// Array          <-    "<array>" Values "</array>" / "<array/>"
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			save := p.ParserData.Pos()
			{
				accept = true
				s := p.ParserData.Pos()
				if p.ParserData.Read() != '<' || p.ParserData.Read() != 'a' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'a' || p.ParserData.Read() != 'y' || p.ParserData.Read() != '>' {
					p.ParserData.Seek(s)
					accept = false
				}
			}
			if accept {
				accept = p.Values()
				if accept {
					{
						accept = true
						s := p.ParserData.Pos()
						if p.ParserData.Read() != '<' || p.ParserData.Read() != '/' || p.ParserData.Read() != 'a' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'a' || p.ParserData.Read() != 'y' || p.ParserData.Read() != '>' {
							p.ParserData.Seek(s)
							accept = false
						}
					}
					if accept {
					}
				}
			}
			if !accept {
				if p.LastError < p.ParserData.Pos() {
					p.LastError = p.ParserData.Pos()
				}
				p.ParserData.Seek(save)
			}
		}
		if !accept {
			{
				accept = true
				s := p.ParserData.Pos()
				if p.ParserData.Read() != '<' || p.ParserData.Read() != 'a' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'r' || p.ParserData.Read() != 'a' || p.ParserData.Read() != 'y' || p.ParserData.Read() != '/' || p.ParserData.Read() != '>' {
					p.ParserData.Seek(s)
					accept = false
				}
			}
			if !accept {
			}
		}
		if !accept {
			p.ParserData.Seek(save)
		}
	}
	end := p.ParserData.Pos()
	if accept {
		node := p.Root.Cleanup(start, end)
		node.Name = "Array"
		node.P = p
		node.Range = node.Range.Clip(p.IgnoreRange)
		p.Root.Append(node)
	} else {
		p.Root.Discard(start)
	}
	if p.IgnoreRange.A >= end || p.IgnoreRange.B <= start {
		p.IgnoreRange = text.Region{}
	}
	return accept
}

func (p *PLIST) Spacing() bool {
	// Spacing        <-    [ \t\n\r]+
	accept := false
	accept = true
	start := p.ParserData.Pos()
	{
		save := p.ParserData.Pos()
		{
			accept = false
			c := p.ParserData.Read()
			if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				accept = true
			} else {
				p.ParserData.UnRead()
			}
		}
		if !accept {
			p.ParserData.Seek(save)
		} else {
			for accept {
				{
					accept = false
					c := p.ParserData.Read()
					if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
						accept = true
					} else {
						p.ParserData.UnRead()
					}
				}
			}
			accept = true
		}
	}
	if accept && start != p.ParserData.Pos() {
		if start < p.IgnoreRange.A || p.IgnoreRange.A == 0 {
			p.IgnoreRange.A = start
		}
		p.IgnoreRange.B = p.ParserData.Pos()
	}
	return accept
}

func (p *PLIST) EndOfFile() bool {
	// EndOfFile      <-    !.
	accept := false
	accept = true
	start := p.ParserData.Pos()
	s := p.ParserData.Pos()
	if p.ParserData.Pos() >= p.ParserData.Len() {
		accept = false
	} else {
		p.ParserData.Read()
		accept = true
	}
	p.ParserData.Seek(s)
	p.Root.Discard(s)
	accept = !accept
	end := p.ParserData.Pos()
	if accept {
		node := p.Root.Cleanup(start, end)
		node.Name = "EndOfFile"
		node.P = p
		node.Range = node.Range.Clip(p.IgnoreRange)
		p.Root.Append(node)
	} else {
		p.Root.Discard(start)
	}
	if p.IgnoreRange.A >= end || p.IgnoreRange.B <= start {
		p.IgnoreRange = text.Region{}
	}
	return accept
}
