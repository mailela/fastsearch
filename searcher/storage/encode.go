package storage

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strings"
)

func Obj2Map(obj interface{}) map[string]interface{} {
	if obj == nil {
		return make(map[string]interface{})
	}
	v := reflect.ValueOf(obj)
	m := make(map[string]interface{}, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i).Name
		value := v.FieldByName(field).Interface()
		m[strings.ToLower(field)] = value
	}
	return m
}

func Encoder(data interface{}) []byte {
	if data == nil {
		return nil
	}
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		// panic(err)
		fmt.Println(err.Error())
	}
	return buffer.Bytes()
}

func Decoder(data []byte, v interface{}) {
	if data == nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(v)
	if err != nil {
		// panic(err)
		fmt.Println(err.Error())
	}
}
