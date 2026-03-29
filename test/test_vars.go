package test

import "os"

// todo: move into main_test.go and add utility functions in there
type TestVars struct {
	DbPath string
}

func NewTestVars() *TestVars {
	return &TestVars{
		DbPath: os.TempDir() + "/dendrite_test_vault",
	}
}


