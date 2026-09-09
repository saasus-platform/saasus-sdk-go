package testdata

import (
	"testing"
)

// TestLoadTestParams は LoadTestParams 関数が正しく動作することを確認します
func TestLoadTestParams(t *testing.T) {
	// LoadTestParams を呼び出す
	params := LoadTestParams(t)

	// 基本的な構造が読み込まれていることを確認
	if params == nil {
		t.Fatal("LoadTestParams returned nil")
	}

	// いくつかのフィールドが存在することを確認
	t.Log("✅ LoadTestParams successfully loaded test parameters")

	// BasicInfo の確認
	if params.BasicInfo.UpdateParams == nil {
		t.Error("BasicInfo.UpdateParams is nil")
	} else {
		t.Logf("BasicInfo.UpdateParams has %d entries", len(params.BasicInfo.UpdateParams))
	}

	// Users の確認
	if params.Users.CreateParams == nil {
		t.Error("Users.CreateParams is nil")
	} else {
		t.Logf("Users.CreateParams has %d entries", len(params.Users.CreateParams))

		// email フィールドの確認
		if email, ok := params.Users.CreateParams["email"]; ok {
			t.Logf("Users.CreateParams[\"email\"] = %v", email)
		} else {
			t.Error("Users.CreateParams[\"email\"] not found")
		}

		// password フィールドの確認
		if password, ok := params.Users.CreateParams["password"]; ok {
			t.Logf("Users.CreateParams[\"password\"] = %v", password)
		} else {
			t.Error("Users.CreateParams[\"password\"] not found")
		}
	}

	// Roles の確認
	if params.Roles.CreateParams == nil {
		t.Error("Roles.CreateParams is nil")
	} else {
		t.Logf("Roles.CreateParams has %d entries", len(params.Roles.CreateParams))
	}

	// Tenants の確認
	if params.Tenants.CreateParams == nil {
		t.Error("Tenants.CreateParams is nil")
	} else {
		t.Logf("Tenants.CreateParams has %d entries", len(params.Tenants.CreateParams))
	}

	// Envs の確認
	if params.Envs.CreateParams == nil {
		t.Error("Envs.CreateParams is nil")
	} else {
		t.Logf("Envs.CreateParams has %d entries", len(params.Envs.CreateParams))
	}
}

// TestPointerHelpers はポインタヘルパー関数が正しく動作することを確認します
func TestPointerHelpers(t *testing.T) {
	// StringPtr のテスト
	str := "test_string"
	strPtr := StringPtr(str)
	if strPtr == nil {
		t.Error("StringPtr returned nil")
	} else if *strPtr != str {
		t.Errorf("StringPtr: expected %s, got %s", str, *strPtr)
	} else {
		t.Logf("✅ StringPtr works correctly: %s", *strPtr)
	}

	// IntPtr のテスト
	num := 42
	numPtr := IntPtr(num)
	if numPtr == nil {
		t.Error("IntPtr returned nil")
	} else if *numPtr != num {
		t.Errorf("IntPtr: expected %d, got %d", num, *numPtr)
	} else {
		t.Logf("✅ IntPtr works correctly: %d", *numPtr)
	}

	// BoolPtr のテスト
	flag := true
	flagPtr := BoolPtr(flag)
	if flagPtr == nil {
		t.Error("BoolPtr returned nil")
	} else if *flagPtr != flag {
		t.Errorf("BoolPtr: expected %v, got %v", flag, *flagPtr)
	} else {
		t.Logf("✅ BoolPtr works correctly: %v", *flagPtr)
	}
}
