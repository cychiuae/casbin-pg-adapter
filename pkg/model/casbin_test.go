package model

import "testing"

func testNewCasbinRuleFromPTypeAndRuleWithLength(t *testing.T, length int) {
	pType := "p"
	rule := []string{"v0", "v1", "v2", "v3", "v4", "v5"}
	casbinRule := NewCasbinRuleFromPTypeAndRule(pType, rule[0:length])
	var wantCasbinRule CasbinRule
	switch length {
	case 1:
		wantCasbinRule = CasbinRule{
			PType: "p",
			V0:    "v0",
		}
	case 2:
		wantCasbinRule = CasbinRule{
			PType: "p",
			V0:    "v0",
			V1:    "v1",
		}
	case 3:
		wantCasbinRule = CasbinRule{
			PType: "p",
			V0:    "v0",
			V1:    "v1",
			V2:    "v2",
		}
	case 4:
		wantCasbinRule = CasbinRule{
			PType: "p",
			V0:    "v0",
			V1:    "v1",
			V2:    "v2",
			V3:    "v3",
		}
	case 5:
		wantCasbinRule = CasbinRule{
			PType: "p",
			V0:    "v0",
			V1:    "v1",
			V2:    "v2",
			V3:    "v2",
			V4:    "v4",
		}
	case 6:
		wantCasbinRule = CasbinRule{
			PType: "p",
			V0:    "v0",
			V1:    "v1",
			V2:    "v2",
			V3:    "v2",
			V4:    "v4",
			V5:    "v5",
		}
	default:
		t.Fatalf("Expected length from 1 to 6")
	}
	if wantCasbinRule != casbinRule {
		t.Errorf(
			"Test NewCasbinRuleFromPTypeAndRule with length %v. Expected %v but got %v",
			length,
			wantCasbinRule,
			casbinRule,
		)
	}
}

func TestNewCasbinRuleFromPTypeAndRule(t *testing.T) {
	for i := 1; i < 5; i++ {
		testNewCasbinRuleFromPTypeAndRuleWithLength(t, i)
	}
}
