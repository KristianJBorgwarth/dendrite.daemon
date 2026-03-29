package test

import "os"

type TestVars struct {
	DbPath string
}

func NewTestVars() *TestVars {
	return &TestVars{
		DbPath: os.TempDir() + "/dendrite_test_vault",
	}
}

