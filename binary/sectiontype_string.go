// Code generated by "stringer -type=SectionType"; DO NOT EDIT.

package binary

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[section_attributes-0]
	_ = x[section_build-1]
	_ = x[section_enums-2]
	_ = x[section_enumValues-3]
	_ = x[section_classes-4]
	_ = x[section_classFunctions-5]
	_ = x[section_classFields-6]
	_ = x[section_functions-7]
	_ = x[section_dynamicCalls-8]
	_ = x[section_registers-9]
	_ = x[section_instructions-10]
	_ = x[section_constants-11]
	_ = x[section_positions-12]
	_ = x[section_files-13]
	_ = x[section_resources-14]
	_ = x[section_sources-15]
	_ = x[section_sourceLines-16]
	_ = x[section_string-17]
	_ = x[section_bytes-18]
	_ = x[section_kInt-19]
	_ = x[section_kFloat-20]
	_ = x[section_kBool-21]
	_ = x[section_kString-22]
	_ = x[section_kNull-23]
	_ = x[section_kUndefined-24]
	_ = x[section_kRune-25]
	_ = x[section_EOF-26]
}

const _SectionType_name = "section_attributessection_buildsection_enumssection_enumValuessection_classessection_classFunctionssection_classFieldssection_functionssection_dynamicCallssection_registerssection_instructionssection_constantssection_positionssection_filessection_resourcessection_sourcessection_sourceLinessection_stringsection_bytessection_kIntsection_kFloatsection_kBoolsection_kStringsection_kNullsection_kUndefinedsection_kRunesection_EOF"

var _SectionType_index = [...]uint16{0, 18, 31, 44, 62, 77, 99, 118, 135, 155, 172, 192, 209, 226, 239, 256, 271, 290, 304, 317, 329, 343, 356, 371, 384, 402, 415, 426}

func (i SectionType) String() string {
	if i < 0 || i >= SectionType(len(_SectionType_index)-1) {
		return "SectionType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SectionType_name[_SectionType_index[i]:_SectionType_index[i+1]]
}
