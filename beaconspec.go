// Package beaconspec provides methods for parsing a Beacon dump file
package beaconspec

import (
	"bufio"
	"net/url"
	"os"
	"strings"
)

// Stores all of the data from the Metadata header in a Beacon file
type BeaconMetadata struct {
	prefix     string
	target     string
	relation   string
	message    string
	annotation string

	description string
	creator     string
	contact     string
	homepage    string
	feed        string
	timestamp   string
	update      string
}

// One line in a Beacon file contains a source, annotation, then a target
type BeaconLine struct {
	Source     string
	Annotation string
	Target     string
}

// ReadMetadata reads a Beacon file, and parse its metadata into a
// BeaconMetadata struct if it is valid
func ReadMetadata(filename string) (BeaconMetadata, error) {
	f, err := os.Open(filename)
	m := BeaconMetadata{}
	if err != nil {
		return m, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "#") {
			break
		}
		// TODO: What if the metadata doesn't match any of the below?
		switch {
		case strings.HasPrefix(line, "#PREFIX:"):
			m.prefix = extractMetadataValue(line)
		case strings.HasPrefix(line, "#TARGET:"):
			m.target = extractMetadataValue(line)
		case strings.HasPrefix(line, "#RELATION:"):
			m.relation = extractMetadataValue(line)
		case strings.HasPrefix(line, "#MESSAGE:"):
			m.message = extractMetadataValue(line)
		case strings.HasPrefix(line, "#ANNOTATION:"):
			m.annotation = extractMetadataValue(line)
		case strings.HasPrefix(line, "#DESCRIPTION:"):
			m.description = extractMetadataValue(line)
		case strings.HasPrefix(line, "#CREATOR:"):
			m.creator = extractMetadataValue(line)
		case strings.HasPrefix(line, "#CONTACT:"):
			m.contact = extractMetadataValue(line)
		case strings.HasPrefix(line, "#HOMEPAGE:"):
			m.homepage = extractMetadataValue(line)
		case strings.HasPrefix(line, "#FEED:"):
			m.feed = extractMetadataValue(line)
		case strings.HasPrefix(line, "#TIMESTAMP:"):
			m.timestamp = extractMetadataValue(line)
		case strings.HasPrefix(line, "#UPDATE:"):
			m.update = extractMetadataValue(line)
		}
	}
	return m, nil
}

func extractMetadataValue(line string) string {
	s := strings.SplitN(line, ":", 2)
	return strings.TrimSpace(s[1])
}

// ParseLine accepts a single line from a Beacon file, and parses it into a
// BeaconLine struct, combining it with the metadata from the file.
func ParseLine(line string, data BeaconMetadata) (BeaconLine, error) {
	b := BeaconLine{}
	s := strings.Split(line, "|")

	if len(s) == 1 {
		b.Source = joinMetaLinks(data.prefix, s[0])
		b.Target = joinMetaLinks(data.target, s[0])
	} else if len(s) == 2 {
		b.Source = joinMetaLinks(data.prefix, s[0])
		// To resolve ambiguity between target and annotation, we
		// need to check if second item is a qualified url
		if isUrl(s[1]) {
			b.Target = s[1]
		} else {
			b.Target = joinMetaLinks(data.target, s[0])
			b.Annotation = s[1]
		}
	} else if len(s) >= 3 {
		b.Source = joinMetaLinks(data.prefix, s[0])
		b.Annotation = s[1]
		b.Target = joinMetaLinks(data.target, s[2])
	}

	return b, nil
}

func isUrl(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// Joins the metadata value with the line value,
// prefixing if it is a url
// suffixing if it is not
func joinMetaLinks(meta, line string) string {
	var link string
	if isUrl(meta) {
		link, _ = url.JoinPath(meta, line)
	} else {
		link, _ = url.JoinPath(line, meta)
	}
	return link
}
