package model

import (
	"strings"
)

const (
	policyLinePrefix = ", "
)

// Filter defines the filtering rules for a FilteredAdapter's policy. Empty values
// are ignored, but all others must match the filter.
type Filter struct {
	P []string
	G []string
}

// CasbinRule is the model for casbin rule
type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

// NewCasbinRuleFromPTypeAndRule returns a CasbinRule from pType and rule
func NewCasbinRuleFromPTypeAndRule(pType string, rule []string) CasbinRule {
	casbinRule := CasbinRule{
		PType: pType,
	}
	ruleLength := len(rule)
	if ruleLength > 0 {
		casbinRule.V0 = rule[0]
	}
	if ruleLength > 1 {
		casbinRule.V1 = rule[1]
	}
	if ruleLength > 2 {
		casbinRule.V2 = rule[2]
	}
	if ruleLength > 3 {
		casbinRule.V3 = rule[3]
	}
	if ruleLength > 4 {
		casbinRule.V4 = rule[4]
	}
	if ruleLength > 5 {
		casbinRule.V5 = rule[5]
	}
	return casbinRule
}

// NewCasbinRuleFromPTypeAndFilter returns a CasbinRule from pType and filter
func NewCasbinRuleFromPTypeAndFilter(pType string, fieldIndex int, fieldValues ...string) CasbinRule {
	casbinRule := CasbinRule{
		PType: pType,
	}

	idx := fieldIndex + len(fieldValues)
	if fieldIndex <= 0 && idx > 0 {
		casbinRule.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && idx > 1 {
		casbinRule.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && idx > 2 {
		casbinRule.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && idx > 3 {
		casbinRule.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && idx > 4 {
		casbinRule.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && idx > 5 {
		casbinRule.V5 = fieldValues[5-fieldIndex]
	}

	return casbinRule
}

// ToPolicyLine map casbinRule to a policy line used in casbin
func (casbinRule CasbinRule) ToPolicyLine() string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(casbinRule.PType)
	if len(casbinRule.V0) > 0 {
		stringBuilder.WriteString(policyLinePrefix)
		stringBuilder.WriteString(casbinRule.V0)
	}
	if len(casbinRule.V1) > 0 {
		stringBuilder.WriteString(policyLinePrefix)
		stringBuilder.WriteString(casbinRule.V1)
	}
	if len(casbinRule.V2) > 0 {
		stringBuilder.WriteString(policyLinePrefix)
		stringBuilder.WriteString(casbinRule.V2)
	}
	if len(casbinRule.V3) > 0 {
		stringBuilder.WriteString(policyLinePrefix)
		stringBuilder.WriteString(casbinRule.V3)
	}
	if len(casbinRule.V4) > 0 {
		stringBuilder.WriteString(policyLinePrefix)
		stringBuilder.WriteString(casbinRule.V4)
	}
	if len(casbinRule.V5) > 0 {
		stringBuilder.WriteString(policyLinePrefix)
		stringBuilder.WriteString(casbinRule.V5)
	}

	return stringBuilder.String()
}

// ToStringSlice map casbinRule to a string slice used in casbin.Model
func (casbinRule CasbinRule) ToStringSlice() []string {
	rule := make([]string, 0)
	rule = append(rule, casbinRule.PType)

	if len(casbinRule.V0) > 0 {
		rule = append(rule, casbinRule.V0)
	}
	if len(casbinRule.V1) > 0 {
		rule = append(rule, casbinRule.V1)
	}
	if len(casbinRule.V2) > 0 {
		rule = append(rule, casbinRule.V2)
	}
	if len(casbinRule.V3) > 0 {
		rule = append(rule, casbinRule.V3)
	}
	if len(casbinRule.V4) > 0 {
		rule = append(rule, casbinRule.V4)
	}
	if len(casbinRule.V5) > 0 {
		rule = append(rule, casbinRule.V5)
	}
	return rule
}
