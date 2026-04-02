package files

import "os"

func WriteToFile(path string, data []byte) (filePath string, err error) {
	if checkIfFileExists(path) {
		return path, nil
	}
	err = os.WriteFile(path, data, 0o644)
	if err != nil {
		return "", err
	}
	return path, nil
}

func checkIfFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
