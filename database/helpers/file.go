package helpers

import (
	"encoding/json"
	"os"
)

func FileExists(f string) bool {
	_, err := os.Stat(f)
	return err == nil
}

// func FileWriteable(f string) bool {
// 	var err error
// 	switch runtime.GOOS {
// 	case "linux":
// 		err = syscall.Access(f, syscall.O_RDWR)
// 	case "windows":
// 		info, _ := os.Stat(f)
// 		m := info.Mode()
// 		return m.IsRegular()

// 	}
// 	return err == nil

// }

func ReadJsonToMap(f string) (map[string]interface{}, error) {
	bytes, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{}, 10)
	_ = json.Unmarshal(bytes, &data)

	return data, nil
}
