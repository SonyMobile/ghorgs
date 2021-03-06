//
// Copyright (c) 2019 Sony Mobile Communications Inc.
// SPDX-License-Identifier: MIT
//

package model

import (
	"fmt"
	"log"
	"sort"
)

type Table struct {
	Records    map[string][]string
	Keys       []string
	Fields     []Field
	pivotField Field
}

func MakeTable(fields []Field) *Table {
	keys := make([]string, 0)
	return &Table{Records: nil, Keys: keys, Fields: fields}
}

func (t *Table) FieldNames() []string {
	return namesOf(t.Fields)
}

func (t *Table) AddKey(key string) {
	t.Keys = append(t.Keys, key)
}

func (t *Table) AddRecord(key string, record []string) {
	if t.Records == nil {
		t.Records = make(map[string][]string)
	}
	t.Records[key] = record
}

func (t *Table) String() string {
	var s string
	for id, key := range t.Keys {
		s += key + "\t"
		line := t.Records[key]
		for i, cell := range line {
			if len(cell) == 0 {
				cell = "-"
			}
			s += cell
			if i < len(line)-1 {
				s += "\t"
			}
		}
		if id < len(t.Keys)-1 {
			s += "\n"
		}
	}

	return s
}

func (t *Table) Log() {
	log.Println(t)
}

// sort interface + method
type By Table

func (a By) Len() int { return len(a.Keys) }
func (a By) Less(i, j int) bool {
	if a.pivotField.Index == ID.Index {
		return a.Keys[i] < a.Keys[j]
	}
	return a.Records[a.Keys[i]][a.pivotField.Index] < a.Records[a.Keys[j]][a.pivotField.Index]
}
func (a By) Swap(i, j int) {
	a.Keys[i], a.Keys[j] = a.Keys[j], a.Keys[i]
}

// sort by field ascending
func (t *Table) SortByField(field string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	tt := *t
	sort.Sort(By(tt))
	return &tt, nil
}

// gets all records with field equals value
func (t *Table) FindAllByField(field, val string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	var ret *Table
	for _, key := range t.Keys {
		if t.Records[key][t.pivotField.Index] == val {
			if ret == nil {
				ret = MakeTable(t.Fields)
			}
			ret.AddKey(key)
			ret.AddRecord(key, t.Records[key])
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("`%s` not found at `%s`", val, field)
	}

	return ret, nil
}

func (t *Table) FindAllByFieldValues(field string, vals []string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	valsmap := make(map[string]int)
	for _, val := range vals {
		valsmap[val] = 0
	}

	var ret *Table
	for _, key := range t.Keys {
		val := t.Records[key][t.pivotField.Index]
		if _, ok := valsmap[val]; ok {
			valsmap[val]++
			if ret == nil {
				ret = MakeTable(t.Fields)
			}
			ret.AddKey(key)
			ret.AddRecord(key, t.Records[key])
		}
	}

	if ret == nil {
		return nil, fmt.Errorf("Requested values not found.")
	}

	ok := true
	for k, v := range valsmap {
		if v == 0 {
			ok = false
			fmt.Printf("`%s` not found in `%s`\n", k, field)
		}
	}

	if !ok {
		return ret, fmt.Errorf("Error!")
	}

	return ret, nil
}

func (t *Table) LessThanByField(field, val string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	ret := MakeTable(t.Fields)
	for _, key := range t.Keys {
		if t.Records[key][t.pivotField.Index] < val {
			ret.AddKey(key)
			ret.AddRecord(key, t.Records[key])
		}
	}

	return ret, nil
}

func (t *Table) GreaterThanByField(field, val string) (*Table, error) {
	err := t.setPivotField(field)
	if err != nil {
		return nil, err
	}

	ret := MakeTable(t.Fields)
	for _, key := range t.Keys {
		if t.Records[key][t.pivotField.Index] > val {
			ret.AddKey(key)
			ret.AddRecord(key, t.Records[key])
		}
	}

	return ret, nil
}

func (t *Table) Last(n int) (*Table, error) {
	if n < 1 || n > len(t.Keys) {
		return nil, fmt.Errorf("Out of range error.")
	}

	keys := make([]string, 0)
	ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
	for _, key := range t.Keys[len(t.Keys)-n:] {
		ret.AddKey(key)
		ret.AddRecord(key, t.Records[key])
	}

	return ret, nil
}

func (t *Table) First(n int) (*Table, error) {
	if n <= 0 || n > len(t.Keys) {
		return nil, fmt.Errorf("Out of range error.")
	}

	keys := make([]string, 0)
	ret := &Table{Records: nil, Keys: keys, Fields: t.Fields}
	for _, key := range t.Keys[:n] {
		ret.AddKey(key)
		ret.AddRecord(key, t.Records[key])
	}

	return ret, nil
}

func (t *Table) setPivotField(fieldName string) error {
	if fieldName == ID.Name {
		t.pivotField = ID
		return nil
	}

	t.pivotField = INVALID_FIELD
	for _, field := range t.Fields {
		if field.Name == fieldName {
			t.pivotField = field
			break
		}
	}

	if t.pivotField == INVALID_FIELD {
		return fmt.Errorf("Invalid search field: %s\n", fieldName)
	}

	return nil
}
