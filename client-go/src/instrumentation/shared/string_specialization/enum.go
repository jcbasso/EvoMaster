package string_specialization

// StringSpecialization Based on taint analysis, could check how input strings are used,
// and inform the search about it
type StringSpecialization int

const (
	// DATE_FORMAT_UNKNOWN_PATTERN String used as a Date with unknown format
	DATE_FORMAT_UNKNOWN_PATTERN StringSpecialization = iota
	// DATE_FORMAT_PATTERN String used as a Date with not explicitly supported format
	DATE_FORMAT_PATTERN
	// DATE_YYYY_MM_DD String used as a Date in YYYY_MM_DD format
	DATE_YYYY_MM_DD
	// DATE_YYYY_MM_DD_HH_MM String used as a Date in YYYY_MM_DD_HH_MM format
	DATE_YYYY_MM_DD_HH_MM
	// ISO_LOCAL_DATE_TIME An ISO Local Date Time (i.e. ISO_LOCAL_DATE + 'T' + ISO_LOCAL_TIME)
	ISO_LOCAL_DATE_TIME
	// ISO_LOCAL_TIME An ISO Local Time (with or without no seconds)
	ISO_LOCAL_TIME
	// INTEGER String used as an integer
	INTEGER
	// CONSTANT String used with a specific, constant value
	CONSTANT
	// CONSTANT_IGNORE_CASE String used with a specific, constant value, ignoring its case
	CONSTANT_IGNORE_CASE
	// REGEX String constrained by a regular expression
	REGEX
	// DOUBLE String parsed to double
	DOUBLE
	// LONG String parsed to long
	LONG
	// BOOLEAN String parsed to boolean
	BOOLEAN
	// FLOAT String parsed to float
	FLOAT
	// EQUAL String should be equal to another string variable,
	// ie 2 (or more) different variables should be keep their
	// value in sync
	EQUAL
)

func (s StringSpecialization) String() string {
	switch s {
	case DATE_FORMAT_UNKNOWN_PATTERN:
		return "DATE_FORMAT_UNKNOWN_PATTERN"
	case DATE_FORMAT_PATTERN:
		return "DATE_FORMAT_PATTERN"
	case DATE_YYYY_MM_DD:
		return "DATE_YYYY_MM_DD"
	case DATE_YYYY_MM_DD_HH_MM:
		return "DATE_YYYY_MM_DD_HH_MM"
	case ISO_LOCAL_DATE_TIME:
		return "ISO_LOCAL_DATE_TIME"
	case ISO_LOCAL_TIME:
		return "ISO_LOCAL_TIME"
	case INTEGER:
		return "INTEGER"
	case CONSTANT:
		return "CONSTANT"
	case CONSTANT_IGNORE_CASE:
		return "CONSTANT_IGNORE_CASE"
	case REGEX:
		return "REGEX"
	case DOUBLE:
		return "DOUBLE"
	case LONG:
		return "LONG"
	case BOOLEAN:
		return "BOOLEAN"
	case FLOAT:
		return "FLOAT"
	case EQUAL:
		return "EQUAL"
	default:
		return ""
	}
}
