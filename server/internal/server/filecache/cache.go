package filecache

import (
	"strings"

	"github.com/takaaa220/swaggo-ide/server/internal/server/protocol"
)

type FileCache struct {
	cache map[protocol.DocumentUri]FileInfo // map[protocol.DocumentUri]FileText
}

type FileInfo struct {
	Version int
	Text    FileText
}

type FileText []string

func (t FileText) GetLine(lineNumber int) (string, bool) {
	if lineNumber < 0 || lineNumber >= len(t) {
		return "", false
	}

	return t[lineNumber], true
}

func (t FileText) String() string {
	return strings.Join(t, "\n")
}

func NewFileText(text string) FileText {
	return FileText(strings.Split(text, "\n"))
}

func (t FileText) Update(changeEvents []protocol.TextDocumentContentChangeEvent) FileText {
	if len(changeEvents) == 0 {
		return t
	}

	return NewFileText(changeEvents[0].Text)
}

func NewFileInfo(version int, text FileText) FileInfo {
	return FileInfo{
		Version: version,
		Text:    text,
	}
}

func NewFileCache() *FileCache {
	return &FileCache{
		cache: make(map[protocol.DocumentUri]FileInfo),
	}
}

func (c *FileCache) Get(uri protocol.DocumentUri) (FileInfo, bool) {
	info, ok := c.cache[uri]
	return info, ok
}

func (c *FileCache) Set(uri protocol.DocumentUri, new FileInfo) {
	current, ok := c.Get(uri)
	if ok && current.Version >= new.Version && new.Version != 0 {
		return
	}

	c.cache[uri] = new
}

func (c *FileCache) Delete(uri protocol.DocumentUri) {
	delete(c.cache, uri)
}
