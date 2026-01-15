package cbor

// Version information for the CBOR library.
const (
	// Version is the current version of the library.
	Version = "1.0.0"

	// VersionMajor is the major version number.
	VersionMajor = 1

	// VersionMinor is the minor version number.
	VersionMinor = 0

	// VersionPatch is the patch version number.
	VersionPatch = 0
)

// VersionInfo returns the full version string.
func VersionInfo() string {
	return "cbor.go v" + Version
}
