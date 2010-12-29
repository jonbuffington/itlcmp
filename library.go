//
// Copyright 2010 by Jon Buffington. All rights reserved.
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
	"container/vector"
	"fmt"
	"http"
	"os"
	"path"
	"strings"
)


const (
	KEY_KIND       = "Kind"       // string
	KEY_LOCATION   = "Location"   // string
	KEY_TRACKS     = "Tracks"     // dict
	KEY_TRACK_TYPE = "Track Type" // string
)


type library struct {
	file     *os.File
	nofiles  vector.StringVector
	notracks vector.StringVector
	stats    *statistics
	tracks   map[string]bool
}


func newLibrary() *library {
	statistics := newStatistics()
	return &library{
		tracks:   make(map[string]bool),
		notracks: vector.StringVector{},
		nofiles:  vector.StringVector{},
		stats:    statistics,
	}
}

func (*library) dir() string {
	return os.Getenv("HOME") + "/Music/iTunes/"
}

func (l *library) open() (err os.Error) {
	l.file, err = os.Open(l.dir()+"iTunes Music Library.xml", os.O_RDONLY, 0)
	return
}

func (l *library) close() {
	l.file.Close()
}

func (l *library) examine() {
	defer func() {
		if ex := recover(); ex == os.EOF {
			// Any missing media files?
			if len(l.nofiles) > 0 {
				fmt.Fprintln(os.Stderr, "\nThe following tracks were not found in your media files:")
				for _, pathname := range l.nofiles {
					fmt.Fprintln(os.Stderr, pathname)
				}
				fmt.Fprintln(os.Stderr)
			}
			// Evaluate media files when EOF is reached.
			fmt.Println("Starting media directory evaluation…")
			// Verify files exist in accumulated library.
			l.checkDir(l.dir() + "iTunes Music")
			// Show the executions stats.
			l.stats.print()
			// Any missing library tracks?
			if len(l.notracks) > 0 {
				fmt.Fprintln(os.Stderr, "\nThe following media files were not found in your tracks:")
				for _, pathname := range l.notracks {
					fmt.Fprintln(os.Stderr, pathname)
				}
			}
		}
	}()

	fmt.Println("Starting library tracks evaluation…")
	// Parse the Library XML file till EOF is reached.
	l.parse()
}

func (l *library) parse() {
	p := newParser(l.file)
	if p.findElem(&plistStartEl, nil) && p.findElem(&dictStartEl, &plistEndEl) && p.findKey(KEY_TRACKS) {
		if p.findElem(&dictStartEl, &dictEndEl) {
			// Loop through the Tracks dictionary.
			for {
				if p.findElem(&keyStartEl, &dictEndEl) && p.findElem(&dictStartEl, &dictEndEl) {
					// TODO: The following strategy is order dependent. Instead, I should be searching
					// for a list of key elements and return which one is located first.
					if p.findKey(KEY_KIND) && p.findElem(&stringStartEl, &dictEndEl) {
						if trackKind := p.text(); len(trackKind) > 0 {
							l.stats.kinds[trackKind]++
						}
					}
					if p.findKey(KEY_LOCATION) && p.findElem(&stringStartEl, &dictEndEl) {
						if location := p.text(); len(location) > 0 {
							if url, err := http.ParseURL(location); err == nil {
								if l.checkTrack(url.Path) {
									l.tracks[url.Path] = true
								}
								l.stats.tracks++
							}
						}
					}
				} else {
					logger.Println("Failed to find Key/Value pair in Tracks dictionary.")
					break
				}
			}
		}
	}
}

func (l *library) checkDir(dirpath string) {
	dir, err := os.Open(dirpath, os.O_RDONLY, 0)
	if err == nil {
		fis, err := dir.Readdir(-1)
		dir.Close()
		if err == nil {
			for _, fi := range fis {
				pathname := dirpath + "/" + fi.Name
				switch {
				case fi.IsRegular():
					l.checkFile(pathname)
				case fi.IsDirectory() && strings.ToLower(path.Ext(pathname)) != ".itlp":
					l.checkDir(pathname)
				}
			}
		}
	}
}

func (l *library) checkFile(pathname string) {
	l.stats.files++
	ext := strings.ToLower(path.Ext(pathname))
	l.stats.exts[ext]++
	if _, ok := l.tracks[pathname]; !ok {
		switch ext {
		case ".epub", ".m4r", ".pdf", ".plist":
			// Ignore non-audio media and ringtones.
		default:
			l.notracks.Push(pathname)
		}
	}
}

func (l *library) checkTrack(pathname string) bool {
	_, err := os.Stat(pathname)
	if err == nil {
		return true
	}
	l.nofiles.Push(pathname)
	return false
}
