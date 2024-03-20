package config

import (
	"errors"
	"fmt"
	"strings"
)

type Value struct {
	Value   string
	Comment bool
}

type Data struct {
	Order []string
	Data  map[string]Value
}

func ParseConfigLine(line string) (string, Value, bool) {
	isComment := strings.HasPrefix(strings.TrimSpace(line), "#")
	if isComment {
		line = strings.TrimPrefix(line, "#")
	}
	pair := strings.Split(line, "|")
	if len(pair) == 2 {
		key := pair[0]
		value := strings.TrimSpace(pair[1])
		return key, Value{Value: value, Comment: isComment}, true
	}
	return "", Value{}, false
}

func ParseConfig(input string) Data {
	lines := strings.Split(input, "\n")
	var config Data
	config.Data = make(map[string]Value)
	for _, line := range lines {
		if key, val, ok := ParseConfigLine(line); ok {
			config.Order = append(config.Order, key)
			config.Data[key] = val
		}
	}
	return config
}

func (d *Data) AddEntry(key string, value Value) (bool, error) {
	// Don't add if the key already exists.
	if _, ok := d.Data[key]; ok {
		return false, errors.New(fmt.Sprintf("The key %s already exists.\n", key))
	}

	d.Order = append(d.Order, key)
	d.Data[key] = value
	return true, nil
}

func (d *Data) RemoveEntry(key string) (bool, error) {
	// Return if the key doesn't exist.
	if _, ok := d.Data[key]; !ok {
		return false, errors.New(fmt.Sprintf("The key %s does not exist.\n", key))
	}

	// Remove the key from the order slice.
	for i, k := range d.Order {
		if k == key {
			d.Order = append(d.Order[:i], d.Order[i+1:]...)
			break
		}
	}

	// Remove the key from the data map.
	delete(d.Data, key)
	return true, nil
}

func (d *Data) ModifyEntry(key string, newValue Value) (bool, error) {
	// Check if the key exists
	if _, ok := d.Data[key]; !ok {
		return false, errors.New(fmt.Sprintf("The key %s does not exist.\n", key))
	}
	// Modify the value for the key
	d.Data[key] = newValue
	return true, nil
}

func (d *Data) GetEntry(key string) (Value, error) {
	if value, ok := d.Data[key]; ok {
		return value, nil
	}
	return Value{}, errors.New(fmt.Sprintf("The key %s does not exist.\n", key))
}

func (d *Data) BuildConfig() string {
	var sb strings.Builder
	for _, key := range d.Order {
		value := d.Data[key]
		if value.Comment {
			sb.WriteString("#")
		}
		sb.WriteString(key + "|" + value.Value + "\n")
	}
	sb.WriteString("RTENDMARKERBS1001")
	return sb.String()
}

func (d *Data) PrintConfig() {
	for _, key := range d.Order {
		val := d.Data[key]
		fmt.Printf("Key:%s, Value:%v\n", key, val)
	}
}
