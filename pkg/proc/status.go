package proc

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type Status struct {
	Name   string  `proc:"Name"`
	Parent Process `proc:"PPid"`
}

func (p *Process) Status() (*Status, error) {
	data, err := p.readFile("status")
	if err != nil {
		return nil, err
	}
	return parseStatus(data)
}

func (p *Process) IsInHierarchyOf(other Process) bool {
	for other != NoProcess {
		if other == *p {
			return true
		}
		stat, err := other.Status()
		if err != nil {
			return false
		}
		other = stat.Parent
	}
	return false
}

func parseStatus(data []byte) (*Status, error) {
	var status Status

	values := make(map[string]string)
	for _, line := range strings.Split(string(data), "\n") {
		key, val, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		val = strings.TrimSpace(val)
		values[key] = val
	}

	v := reflect.ValueOf(&status)

	t := v.Elem().Type()
	for i := 0; i < t.NumField(); i++ {
		fv := t.Field(i)
		tags := strings.Split(fv.Tag.Get("proc"), ",")
		tagName := fv.Name
		if len(tags) > 0 {
			tagName = tags[0]
		}
		value, ok := values[tagName]
		if !ok {
			continue
		}
		subject := v.Elem().Field(i)

		if !v.Elem().CanSet() {
			return nil, fmt.Errorf("target is not settable")
		}

		switch subject.Kind() {
		case reflect.String:
			subject.SetString(value)
		case reflect.Uint64:
			u, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			subject.SetUint(uint64(u))
		default:
			return nil, fmt.Errorf("decoding of kind %s is not supported", subject.Kind())
		}
	}

	return &status, nil
}
