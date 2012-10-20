//
// Copyright 2010-2012 by Jon Buffington. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
  "os"
	"xml"
)


var (
	dictStartEl   = xml.StartElement{Name: xml.Name{Local: "dict"}}
	dictEndEl     = xml.EndElement{Name: xml.Name{Local: "dict"}}
	keyStartEl    = xml.StartElement{Name: xml.Name{Local: "key"}}
	keyEndEl      = xml.EndElement{Name: xml.Name{Local: "key"}}
	plistStartEl  = xml.StartElement{Name: xml.Name{Local: "plist"}}
	plistEndEl    = xml.EndElement{Name: xml.Name{Local: "plist"}}
	stringStartEl = xml.StartElement{Name: xml.Name{Local: "string"}}
	stringEndEl   = xml.EndElement{Name: xml.Name{Local: "string"}}
)


type parser xml.Parser


func newParser(file *os.File) *parser {
  return (*parser)(xml.NewParser(file))
}

func (p *parser) nextToken() xml.Token {
	tok, err := (*xml.Parser)(p).Token()
	if err != nil {
		// Bounce up the call train to the recover handler.
		panic(err)
	}
	return tok
}

func (p *parser) findElem(target *xml.StartElement, until *xml.EndElement) bool {
	for {
		tok := p.nextToken()
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == target.Name.Local {
				return true
			}
		case xml.EndElement:
			if t.Name.Local == until.Name.Local {
				break
			}
		}
	}
	return false
}

func (p *parser) findKey(keyname string) bool {
	for {
		if p.findElem(&keyStartEl, &dictEndEl) {
			// Locate the character data that matches the keyname.
			if value := p.text(); value == keyname {
				return true
			}
		}
	}
	return false
}

func (p *parser) text() string {
	tok := p.nextToken()
	if t, isCharData := tok.(xml.CharData); isCharData {
		return string([]byte(t))
	}
	return ""
}
