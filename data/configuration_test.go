package data

import (
	"testing"
)

func Test_Device_Success(t *testing.T) {
	d := Configuration{Code: "test12345"}
	m1 := Module{Code: "mod1"}
	if mid, err := d.Module(m1); err != nil {
		t.Errorf("error create module %v", err)
	} else {
		if mid != 0 {
			t.Error("mid != 0")
		}
		p11 := Property{Type: 2, Code: "prop11"}
		if pid, err := d.Property(m1.Code, p11); err != nil {
			t.Errorf("error create property %v", err)
		} else if pid != 0 {
			t.Error("property id != 0")
		}
		p12 := Property{Type: 3, Code: ("prop12")}
		if pid, err := d.Property(m1.Code, p12); err != nil {
			t.Errorf("error create property %v", err)
		} else if pid != 1 {
			t.Error("property id != 1")
		}
	}
	m2 := Module{Code: "mod2"}
	if mid, err := d.Module(m2); err != nil {
		t.Errorf("error create module %v", err)
	} else {
		if mid != 1 {
			t.Error("mid != 0")
		}
		p21 := Property{Type: 2, Code: ("prop21")}
		if pid, err := d.Property(m2.Code, p21); err != nil {
			t.Errorf("error create property %v", err)
		} else if pid != 2 {
			t.Error("property id != 2")
		}
		p22 := Property{Type: 3, Code: ("prop22")}
		if pid, err := d.Property(m2.Code, p22); err != nil {
			t.Errorf("error create property %v", err)
		} else if pid != 3 {
			t.Error("property id != 3")
		}
	}
}
