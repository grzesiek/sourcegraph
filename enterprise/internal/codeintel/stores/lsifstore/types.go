package lsifstore

import "github.com/sourcegraph/sourcegraph/lib/codeintel/precise"

// Package pairs a package name and the dump that provides it.
type Package struct {
	DumpID  int
	Scheme  string
	Name    string
	Version string
}

// PackageReferences pairs a package name/version with a dump that depends on it.
type PackageReference struct {
	Package
	Filter []byte // a bloom filter of identifiers imported by this dependent
}

// Location is an LSP-like location scoped to a dump.
type Location struct {
	DumpID int
	Path   string
	Range  Range
}

// Range is an inclusive bounds within a file.
type Range struct {
	Start Position
	End   Position
}

// Position is a unique position within a file.
type Position struct {
	Line      int
	Character int
}

// Diagnostic describes diagnostic information attached to a location within a
// particular dump.
type Diagnostic struct {
	DumpID int
	Path   string
	precise.DiagnosticData
}

// CodeIntelligenceRange pairs a range with its definitions, reference, hover text, and documentation.
type CodeIntelligenceRange struct {
	Range               Range
	Definitions         []Location
	References          []Location
	HoverText           string
	DocumentationPathID string
}
