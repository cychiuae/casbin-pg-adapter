package model

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
