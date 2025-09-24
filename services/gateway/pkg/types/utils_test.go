package types

import "encoding/json"

func toJson(obj any) []byte {
	b, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return b
}

func unmarshalFromRawJson[T any](rawJson []byte) (b T, err error) {
	err = json.Unmarshal(rawJson, &b)
	return
}

func unmarshalFromString[T any](s string) (T, error) {
	asJson := toJson(s)
	return unmarshalFromRawJson[T](asJson)
}
