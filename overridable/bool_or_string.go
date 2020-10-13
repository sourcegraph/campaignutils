package overridable

import (
	"encoding/json"
)

// BoolOrString is a set of rules that either evaluate to a string or a bool.
type BoolOrString struct {
	rules rules
}

// FromBoolOrString creates a BoolOrString representing a static, scalar value.
func FromBoolOrString(v interface{}) BoolOrString {
	return BoolOrString{
		rules: rules{SimpleRule(v)},
	}
}

// Value returns the value for the given repository.
func (bs *BoolOrString) Value(name string) interface{} {
	v := bs.rules.Match(name)
	if v == nil {
		return false
	}
	return v
}

// MarshalJSON encodes the BoolOrString overridable to a json representation.
func (bs *BoolOrString) MarshalJSON() ([]byte, error) {
	if len(bs.rules) == 0 {
		return []byte("false"), nil
	}
	return json.Marshal(bs.rules)
}

// UnmarshalJSON unmarshalls a JSON value into a Publish.
func (bs *BoolOrString) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		*bs = BoolOrString{rules: rules{SimpleRule(b)}}
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*bs = BoolOrString{rules: rules{SimpleRule(s)}}
		return nil
	}

	var c complex
	if err := json.Unmarshal(data, &c); err != nil {
		return err
	}

	return bs.rules.hydrateFromComplex(c)
}

// UnmarshalYAML unmarshalls a YAML value into a Publish.
func (bs *BoolOrString) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var b bool
	if err := unmarshal(&b); err == nil {
		*bs = BoolOrString{rules: rules{SimpleRule(b)}}
		return nil
	}

	var s string
	if err := unmarshal(&s); err == nil {
		*bs = BoolOrString{rules: rules{SimpleRule(s)}}
		return nil
	}

	var c complex
	if err := unmarshal(&c); err != nil {
		return err
	}

	return bs.rules.hydrateFromComplex(c)
}

// Equal tests two BoolOrStrings for equality, used in cmp.
func (bs *BoolOrString) Equal(other *BoolOrString) bool {
	return bs.rules.Equal(other.rules)
}
