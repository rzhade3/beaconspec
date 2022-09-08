package beaconspec

import (
	"testing"
)

func TestReadMetadata(t *testing.T) {
	actual, err := ReadMetadata("test_beacon_file.txt")
	if err != nil {
		t.Errorf("ReadMetadata() != nil")
	}
	expected := BeaconMetadata{
		prefix:      "http://zhade.dev/",
		target:      "http://zhade.dev/beacon",
		relation:    "http://zhade.dev/relation",
		message:     "This is a message!",
		annotation:  "This is an annotation!",
		description: "This is a description of the filedump!",
		creator:     "Ronald McDonald",
		contact:     "foo@bar.com",
		homepage:    "http://example.com/",
		feed:        "http://example.com/download",
		timestamp:   "2022-09-28T00:00:00.000000",
		update:      "weekly",
	}

	if actual != expected {
		t.Errorf("ParseMetadata() = %+v, actual %+v", expected, actual)
	}
}

func TestExtractMetadataValue(t *testing.T) {
	type args struct {
		line     string
		expected string
	}

	tests := []args{
		{
			line:     "#ARBITRARY: http://example.com/",
			expected: "http://example.com/",
		},
		{
			line:     "#FOO: This can contain any characters$%%",
			expected: "This can contain any characters$%%",
		},
		{
			line:     "#BAR: Excess whitespace is removed  \n\t",
			expected: "Excess whitespace is removed",
		},
	}
	for _, tt := range tests {
		if got := extractMetadataValue(tt.line); got != tt.expected {
			t.Errorf("extractMetadataValue() = %+v, want %+v", got, tt.expected)
		}
	}
}

func TestIsUrl(t *testing.T) {
	s := "http://example.com/"
	if !isUrl(s) {
		t.Errorf("isUrl(%s) != true", s)
	}
	s = "htt/not-url"
	if isUrl(s) {
		t.Errorf("isUrl(%s) != false", s)
	}
}

func TestJoinMetaLinks(t *testing.T) {
	type args struct {
		meta     string
		line     string
		expected string
	}
	tests := []args{
		{
			meta:     "http://example.com/",
			line:     "foo bar",
			expected: "http://example.com/foo%20bar",
		},
		{
			meta:     "bar foo",
			line:     "http://example.com/",
			expected: "http://example.com/bar%20foo",
		},
	}
	for _, tt := range tests {
		if got := joinMetaLinks(tt.meta, tt.line); got != tt.expected {
			t.Errorf("joinMetaLinks() = %+v, want %+v", got, tt.expected)
		}
	}
}

func TestParseLine(t *testing.T) {
	type args struct {
		line     string
		expected BeaconLine
	}

	tests := []args{
		{
			line: "google.com",
			expected: BeaconLine{
				Source:     "google.com",
				Target:     "google.com",
				Annotation: "",
			},
		},
		{
			line: "google.com|foo",
			expected: BeaconLine{
				Source:     "google.com",
				Target:     "google.com",
				Annotation: "foo",
			},
		},
		{
			line: "google.com|https://foo.com",
			expected: BeaconLine{
				Source:     "google.com",
				Target:     "https://foo.com",
				Annotation: "",
			},
		},
		{
			line: "http://source.com/|Should be an annotation|http://target.com/",
			expected: BeaconLine{
				Source:     "http://source.com/",
				Annotation: "Should be an annotation",
				Target:     "http://target.com/",
			},
		},
		{
			line: "http://source.com/||ibn:weirdlookingtarget",
			expected: BeaconLine{
				Source:     "http://source.com/",
				Annotation: "",
				Target:     "ibn:weirdlookingtarget",
			},
		},
	}
	for _, tt := range tests {
		got, err := ParseLine(tt.line, &BeaconMetadata{})
		if err != nil {
			t.Errorf("ParseLine(%s) error %v", tt.line, err)
		}
		if got != tt.expected {
			t.Errorf("ParseLine(%s) = %v, want %v", tt.line, got, tt.expected)
		}
	}
}

func TestParseLineWithMetadata(t *testing.T) {
	type args struct {
		line     string
		meta     BeaconMetadata
		expected BeaconLine
	}

	tests := []args{
		{
			line: "google.com",
			meta: BeaconMetadata{
				prefix: "http://example1.com/",
				target: "http://example2.com",
			},
			expected: BeaconLine{
				Source:     "http://example1.com/google.com",
				Target:     "http://example2.com/google.com",
				Annotation: "",
			},
		},
		{
			line: "google.com|foo",
			meta: BeaconMetadata{
				prefix: "http://example1.com/",
				target: "http://example2.com",
			},
			expected: BeaconLine{
				Source:     "http://example1.com/google.com",
				Target:     "http://example2.com/google.com",
				Annotation: "foo",
			},
		},
	}

	for _, tt := range tests {
		got, err := ParseLine(tt.line, &tt.meta)
		if err != nil {
			t.Errorf("ParseLine(%s) error %v", tt.line, err)
		}
		if got != tt.expected {
			t.Errorf("ParseLine(%s) = %+v, want %+v", tt.line, got, tt.expected)
		}
	}
}

func TestParseLineError(t *testing.T) {
	s := "http://example.com/|foo|bar|baz"
	_, err := ParseLine(s, &BeaconMetadata{})
	if err == nil {
		t.Errorf("ParseLine(%s) error = nil, want error", s)
	}
}
