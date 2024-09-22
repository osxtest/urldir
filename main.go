package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	var depth int
	flag.IntVar(&depth, "d", -1, "Depth of dirs to print; -1 to print all dirs.")
	flag.Parse()

	res := make(map[string]struct{}, 4096)
	size := 2048 * 2048

	sc := bufio.NewScanner(os.Stdin)
	sc.Buffer(make([]byte, size), size)
	for sc.Scan() {
		u := sc.Text()
		for _, v := range parse(u, depth) {
			res[v] = struct{}{}
		}
	}
	for v := range res {
		fmt.Println(v)
	}
}

func parse(u string, depth int) []string {
	if u == "" {
		return nil
	}
	parsed, err := url.Parse(u)
	if err != nil {
		return nil
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil
	}
	if parsed.Host == "" {
		return nil
	}
	return dirs(parsed, depth)
}

// build builds the url with the path.
func build(u *url.URL, path string) string {
	if path == "" {
		return u.Scheme + "://" + u.Host
	}
	if strings.HasPrefix(path, "/") {
		return u.Scheme + "://" + u.Host + path
	}
	return u.Scheme + "://" + u.Host + "/" + path
}

// dirs returns the directories of the url up to the specified depth. If depth
// is -1, it returns all the directories.
// E.g. URL : Depth
// http://example.com : 0
// http://example.com/a	: 1
// http://example.com/a/b : 2
func dirs(u *url.URL, depth int) []string {
	p := path.Clean(u.Path)
	dir := path.Dir(p)

	// Handle path cases like "/a", "/", "".
	if dir == "/" || dir == "." || dir == "" || depth == 0 {
		return []string{build(u, "")}
	}

	// The path has at least one directory, e.g. "/a/b" -> "/a".
	parts := strings.Split(dir, "/")
	res := make([]string, 0, len(parts))
	for i := range parts {
		if depth > 0 && i > depth {
			break
		}
		p := strings.Join(parts[:i+1], "/")
		res = append(res, build(u, p))
	}
	return res
}
