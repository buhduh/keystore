package main

import (
	"testing"
)

func TestSecurity(t *testing.T) {
	t.Run("2Way", twoWay)
}

func twoWay(t *testing.T) {
	salt := "some sweet sweet salt &&&&."
	key := "A super super secret key.!!!"
	phrases := []struct {
		phrase string
	}{
		{"I want to encrypt this."},
		{"How about some more with &^^^^^^^*()*&%*$$)((&(&//\\"},
		{`How about some raw strings?*&^(%asdf cadsf*(&^(2345)(&*%$&\\\asdf\`},
	}
	enc := &TwoWayDecrypt{key: key}
	for _, p := range phrases {
		encrypted, err := enc.EncryptPassword(salt, p.phrase)
		if err != nil {
			t.Logf("Should not error, got '%s'.", err)
			t.Fail()
		}
		decrypted, err := enc.DecryptPassword(salt, encrypted)
		if err != nil {
			t.Logf("Should not error, got '%s'.", err)
			t.Fail()
		}
		if decrypted != p.phrase {
			t.Logf("Expected '%s', got '%s'.", p.phrase, decrypted)
			t.Fail()
		}
	}
}
