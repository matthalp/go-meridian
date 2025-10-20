// Package main implements a code generator for timezone packages.
// It reads timezone definitions from timezones.yaml and generates
// package files and tests for each timezone.
package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config represents the timezones.yaml structure.
type Config struct {
	Timezones []TimezoneDef `yaml:"timezones"`
}

// TimezoneDef defines a single timezone.
type TimezoneDef struct {
	Name        string `yaml:"name"`
	Location    string `yaml:"location"`
	Description string `yaml:"description"`
}

// TemplateData contains all variables needed for template rendering.
type TemplateData struct {
	PackageName string
	Location    string
	Description string
	Abbrev      string
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Println("✓ Successfully generated all timezone packages")
}

func run() error {
	// Read timezones.yaml
	data, err := os.ReadFile("timezones.yaml")
	if err != nil {
		return fmt.Errorf("failed to read timezones.yaml: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse timezones.yaml: %w", err)
	}

	// Generate each timezone package
	for _, tz := range config.Timezones {
		if err := generateTimezone(tz); err != nil {
			return fmt.Errorf("failed to generate %s: %w", tz.Name, err)
		}
		fmt.Printf("Generated %s package\n", tz.Name)
	}

	return nil
}

func generateTimezone(def TimezoneDef) error {
	// Prepare template data
	data := TemplateData{
		PackageName: def.Name,
		Location:    def.Location,
		Description: def.Description,
		Abbrev:      strings.ToUpper(def.Name),
	}

	// Create package directory
	pkgDir := def.Name
	if err := os.MkdirAll(pkgDir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", pkgDir, err)
	}

	// Generate package file
	pkgFile := filepath.Join(pkgDir, def.Name+".go")
	if err := generateFile(pkgFile, packageTemplate, data); err != nil {
		return fmt.Errorf("failed to generate package file: %w", err)
	}

	// Generate test file
	testFile := filepath.Join(pkgDir, def.Name+"_test.go")
	if err := generateFile(testFile, testTemplate, data); err != nil {
		return fmt.Errorf("failed to generate test file: %w", err)
	}

	return nil
}

func generateFile(filename string, tmpl *template.Template, data TemplateData) error {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// Write to file first
	if err := os.WriteFile(filename, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	// Format using goimports (handles both formatting and import ordering)
	cmd := exec.Command("goimports", "-w", filename)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to format with goimports: %w\nOutput: %s", err, output)
	}

	return nil
}

var packageTemplate = template.Must(template.New("package").Parse(`/*
Package {{.PackageName}} provides {{.Description}} timezone support for meridian.
{{if eq .PackageName "utc"}}
{{.Abbrev}} ({{.Description}}) is the primary time standard by which the world
regulates clocks and time. It is timezone-neutral and does not observe daylight
saving time.
{{else}}
{{.Abbrev}} represents the {{.Location}} IANA timezone, which observes {{.Description}}{{if eq .PackageName "est"}} (EST) and Eastern Daylight Time (EDT){{else if eq .PackageName "pst"}} (PST) and Pacific Daylight Time (PDT){{end}} depending on the time of year.
{{end}}
# Usage

Create {{.Abbrev}} times:

	now := {{.PackageName}}.Now()
{{if eq .PackageName "est"}}	meeting := {{.PackageName}}.Date(2024, time.December, 25, 9, 0, 0, 0)
{{else if eq .PackageName "pst"}}	event := {{.PackageName}}.Date(2024, time.December, 25, 6, 0, 0, 0)
{{else}}	specific := {{.PackageName}}.Date(2024, time.December, 25, 10, 30, 0, 0)
{{end}}	parsed, _ := {{.PackageName}}.Parse(time.RFC3339, "2024-12-25T{{if eq .PackageName "est"}}09:00:00-05:00{{else if eq .PackageName "pst"}}06:00:00-08:00{{else}}10:30:00Z{{end}}")
{{if eq .PackageName "utc"}}
Convert to {{.Abbrev}} from other timezones:

	eastern := est.Now()
	universal := {{.PackageName}}.FromMoment(eastern)

The {{.PackageName}}.Time type is an alias for meridian.Time[{{.PackageName}}.Timezone], providing
compile-time timezone safety while maintaining compatibility with standard
time.Time through the Moment interface.
{{else}}
Convert to {{.Abbrev}} from other timezones:

{{if eq .PackageName "est"}}	pacific := pst.Now()
	eastern := {{.PackageName}}.FromMoment(pacific)
{{else}}	eastern := est.Now()
	pacific := {{.PackageName}}.FromMoment(eastern)
{{end}}
Convert from standard time.Time:

	stdTime := time.Now()
	typedTime := {{.PackageName}}.FromMoment(stdTime)

The {{.PackageName}}.Time type is an alias for meridian.Time[{{.PackageName}}.Timezone], providing
compile-time timezone safety. Functions that accept {{.PackageName}}.Time can only receive
times explicitly typed as {{.Description}}, preventing timezone confusion.
{{end}}*/
package {{.PackageName}}

import (
	"fmt"
	"time"

	"github.com/matthalp/go-meridian"
)

// location is the IANA timezone location, loaded once at package initialization.
var location = mustLoadLocation("{{.Location}}")

// mustLoadLocation loads a timezone location or panics if it fails.
// This should only fail if the system's timezone database is corrupted or missing.
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(fmt.Sprintf("failed to load timezone %s: %v", name, err))
	}
	return loc
}

// Timezone represents the {{.Description}} timezone.
type Timezone struct{}

// Location returns the IANA timezone location.
func (Timezone) Location() *time.Location {
	return location
}

// Time is a convenience alias for meridian.Time[Timezone].
type Time = meridian.Time[Timezone]

// Now returns the current time in this timezone.
func Now() Time {
	return meridian.Now[Timezone]()
}

// Date creates a new time in this timezone with the specified date and time components.
func Date(year int, month time.Month, day, hour, minute, sec, nsec int) Time {
	return meridian.Date[Timezone](year, month, day, hour, minute, sec, nsec)
}

// FromMoment converts any Moment to {{.Abbrev}} time.
func FromMoment(m meridian.Moment) Time {
	return meridian.FromMoment[Timezone](m)
}

// Parse parses a formatted string and returns the time value it represents in {{.Abbrev}}.
// The layout defines the format by showing how the reference time would be displayed.
// The time is parsed in the {{.Location}} location.
func Parse(layout, value string) (Time, error) {
	return meridian.Parse[Timezone](layout, value)
}

// Unix returns the {{.Abbrev}} time corresponding to the given Unix time,
// sec seconds and nsec nanoseconds since January 1, 1970 UTC.
func Unix(sec, nsec int64) Time {
	return meridian.Unix[Timezone](sec, nsec)
}

// UnixMilli returns the {{.Abbrev}} time corresponding to the given Unix time,
// msec milliseconds since January 1, 1970 UTC.
func UnixMilli(msec int64) Time {
	return meridian.UnixMilli[Timezone](msec)
}

// UnixMicro returns the {{.Abbrev}} time corresponding to the given Unix time,
// usec microseconds since January 1, 1970 UTC.
func UnixMicro(usec int64) Time {
	return meridian.UnixMicro[Timezone](usec)
}
`))

var testTemplate = template.Must(template.New("test").Parse(`package {{.PackageName}}

import (
	"testing"
	"time"
{{- if or (ne .PackageName "pt") (ne .PackageName "utc")}}

{{- if ne .PackageName "pt"}}
	"github.com/matthalp/go-meridian/pt"
{{- end}}
{{- if ne .PackageName "utc"}}
	"github.com/matthalp/go-meridian/utc"
{{- end}}
{{- end}}
)

func Test{{.Abbrev}}Location(t *testing.T) {
	var tz Timezone
	loc := tz.Location()
	if loc.String() != "{{.Location}}" {
		t.Errorf("Timezone.Location() = %v, want {{.Location}}", loc.String())
	}
}

func TestNow(t *testing.T) {
	before := time.Now().UTC()
	tzTime := Now()
	after := time.Now().UTC()

	// Parse back to verify it's within range
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}

	if parsed.Before(before.Add(-time.Second)) || parsed.After(after.Add(time.Second)) {
		t.Errorf("Now() returned time outside expected range: got %v, expected between %v and %v", parsed, before, after)
	}
}

func TestDate(t *testing.T) {
	// Create a time: Jan 15, 2024 at noon {{.Abbrev}}
	tzTime := Date(2024, time.January, 15, 12, 0, 0, 0)

	// Format should show the time in {{.Abbrev}}
	result := tzTime.Format("15:04 MST")

	// January 15 is during winter, so should show standard time abbreviation
	// The IANA database provides timezone-specific abbreviations (EST, PST, etc.)
	// We just verify it contains the expected hour
	if !contains(result, "12:00") {
		t.Errorf("Format() = %q, expected to contain 12:00", result)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s[:len(substr)] == substr || contains(s[1:], substr))
}

func TestDateWithOffset(t *testing.T) {
	// Create a time in {{.Abbrev}} (UTC offset varies by timezone and DST)
	// Noon {{.Abbrev}} should have corresponding UTC offset
	tzTime := Date(2024, time.January, 1, 12, 0, 0, 0)

	// Parse the formatted time and convert to UTC to verify
	parsed, err := time.Parse(time.RFC3339, tzTime.Format(time.RFC3339))
	if err != nil {
		t.Fatalf("Failed to parse result: %v", err)
	}
	utcTime := parsed.UTC()

	// Verify that the hour in {{.Abbrev}} location is 12
	locationTime := utcTime.In(location)
	if locationTime.Hour() != 12 {
		t.Errorf("Date() hour in {{.Abbrev}} = %v, want 12", locationTime.Hour())
	}
}

func TestFromMoment(t *testing.T) {
	t.Run("from time.Time", func(t *testing.T) {
		// Test converting from standard time.Time in UTC
		stdTime := time.Date(2024, time.January, 15, 17, 0, 0, 0, time.UTC)
		{{.PackageName}}Time := FromMoment(stdTime)

		// Verify the conversion - should represent same moment
		if !{{.PackageName}}Time.UTC().Equal(stdTime) {
			t.Errorf("FromMoment(time.Time) UTC = %v, want %v", {{.PackageName}}Time.UTC(), stdTime)
		}
	})

{{- if ne .PackageName "utc"}}

	t.Run("from UTC", func(t *testing.T) {
		// Create 17:00 UTC
		utcTime := utc.Date(2024, time.January, 15, 17, 0, 0, 0)

		// Convert to {{.Abbrev}}
		{{.PackageName}}Time := FromMoment(utcTime)

		// Verify same moment in time
		if !{{.PackageName}}Time.UTC().Equal(utcTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})
{{- end}}
{{- if ne .PackageName "pt"}}

	t.Run("from PT", func(t *testing.T) {
		// Create 9:00 PT
		ptTime := pt.Date(2024, time.January, 15, 9, 0, 0, 0)

		// Convert to {{.Abbrev}}
		{{.PackageName}}Time := FromMoment(ptTime)

		// Verify same moment in time
		if !{{.PackageName}}Time.UTC().Equal(ptTime.UTC()) {
			t.Error("Converted time doesn't represent same moment")
		}
	})
{{- end}}

{{- if ne .PackageName "utc"}}

	t.Run("round trip conversion", func(t *testing.T) {
		// Create time in {{.Abbrev}}
		original := Date(2024, time.January, 15, 14, 30, 0, 0)

		// Convert to UTC and back
		viaUTC := FromMoment(utc.FromMoment(original))

		// Should represent the same moment
		if !viaUTC.UTC().Equal(original.UTC()) {
			t.Error("Round trip conversion changed the moment in time")
		}

		// Should format the same
		if viaUTC.Format(time.RFC3339) != original.Format(time.RFC3339) {
			t.Errorf("Round trip format = %q, want %q",
				viaUTC.Format(time.RFC3339), original.Format(time.RFC3339))
		}
	})
{{- end}}
}

func TestParse(t *testing.T) {
	t.Run("RFC3339 format", func(t *testing.T) {
		// Parse a time string without timezone, should be interpreted as {{.Abbrev}}
		parsed, err := Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Should be interpreted as 12:00 {{.Abbrev}}
		expected := Date(2024, time.January, 15, 12, 0, 0, 0)
		if parsed.Format(time.RFC3339) != expected.Format(time.RFC3339) {
			t.Errorf("Parse() = %v, want %v", parsed.Format(time.RFC3339), expected.Format(time.RFC3339))
		}
	})

{{- if ne .PackageName "utc"}}

	t.Run("timezone specific interpretation", func(t *testing.T) {
		// Parse same clock time in {{.Abbrev}}
		{{.PackageName}}Parsed, err := Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}

		// Same clock time parsed in UTC would be different
		utcParsed, err := utc.Parse("2006-01-02 15:04:05", "2024-01-15 12:00:00")
		if err != nil {
			t.Fatalf("utc.Parse() error = %v", err)
		}

		// They should represent different moments in time
		if {{.PackageName}}Parsed.UTC().Equal(utcParsed.UTC()) {
			t.Error("{{.Abbrev}} and UTC parse of same clock time should be different moments")
		}
	})
{{- end}}

	t.Run("invalid format", func(t *testing.T) {
		_, err := Parse(time.RFC3339, "invalid-time-string")
		if err == nil {
			t.Error("Parse() expected error for invalid input, got nil")
		}
	})
}

func TestUnix(t *testing.T) {
	t.Run("epoch", func(t *testing.T) {
		epoch := Unix(0, 0)
		
		// But UTC should be epoch
		if !epoch.UTC().Equal(time.Unix(0, 0)) {
			t.Error("Unix(0, 0) UTC time should be epoch")
		}
	})

	t.Run("known timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00 UTC
		result := Unix(1705320000, 0)
		
		// Verify UTC equivalence
		if !result.UTC().Equal(time.Unix(1705320000, 0)) {
			t.Error("Unix timestamp doesn't match")
		}
	})
}

func TestUnixMilli(t *testing.T) {
	t.Run("known millisecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000 UTC
		msec := int64(1705320000000)
		result := UnixMilli(msec)
		
		// Verify UTC equivalence
		if !result.UTC().Equal(time.UnixMilli(msec)) {
			t.Error("UnixMilli UTC time doesn't match")
		}
	})

	t.Run("with milliseconds precision", func(t *testing.T) {
		msec := int64(1705320000123)
		result := UnixMilli(msec)
		if !result.UTC().Equal(time.UnixMilli(msec)) {
			t.Errorf("UnixMilli precision mismatch")
		}
	})
}

func TestUnixMicro(t *testing.T) {
	t.Run("known microsecond timestamp", func(t *testing.T) {
		// 2024-01-15 12:00:00.000000 UTC
		usec := int64(1705320000000000)
		result := UnixMicro(usec)
		
		// Verify UTC equivalence
		if !result.UTC().Equal(time.UnixMicro(usec)) {
			t.Error("UnixMicro UTC time doesn't match")
		}
	})

	t.Run("with microseconds precision", func(t *testing.T) {
		usec := int64(1705320000123456)
		result := UnixMicro(usec)
		if !result.UTC().Equal(time.UnixMicro(usec)) {
			t.Errorf("UnixMicro precision mismatch")
		}
	})
}
`))
