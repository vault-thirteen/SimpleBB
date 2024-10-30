package nt

import (
	"fmt"
	"reflect"
)

const (
	Component_F = "F"
	Component_M = "M"
	Component_R = "R"
	Component_T = "T"
	Component_U = "U"
)

const (
	ErrF_InvalidComponentName  = "invalid component name: %v"
	ErrF_InvalidComponentValue = "invalid value of a component: %v"
)

// Component is a notification template component.
type Component struct {
	// Name is a short component name.
	Name string

	// Value is the value of a component.
	Value any
}

func NewComponent(name string, value any) (c *Component, err error) {
	err = checkComponentName(name)
	if err != nil {
		return nil, err
	}

	err = checkComponentValue(name, value)
	if err != nil {
		return nil, err
	}

	c = &Component{
		Name:  name,
		Value: value,
	}
	return c, nil
}

// checkComponentName validates the component's name.
func checkComponentName(name string) (err error) {
	switch name {
	case Component_F,
		Component_M,
		Component_R,
		Component_T,
		Component_U:
		return nil

	default:
		return fmt.Errorf(ErrF_InvalidComponentName, name)
	}
}

// checkComponentValue validates the component's type according to its name.
func checkComponentValue(name string, value any) (err error) {
	if name == Component_F {
		if reflect.TypeOf(value).Name() != "string" {
			return fmt.Errorf(ErrF_InvalidComponentValue, name)
		} else {
			return nil
		}
	}

	switch name {
	case Component_M,
		Component_R,
		Component_T,
		Component_U:
		if reflect.TypeOf(value).Name() != "int" {
			return fmt.Errorf(ErrF_InvalidComponentValue, name)
		} else {
			return nil
		}

	default:
		return fmt.Errorf(ErrF_InvalidComponentName, name)
	}
}
