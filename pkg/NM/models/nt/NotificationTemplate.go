package nt

import (
	"errors"
	"fmt"

	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

const (
	Err_Stupid                       = "you are stupid"
	Err_ArgsCount                    = "arguments count error"
	Err_PlaceholdersCount            = "placeholders count error"
	Err_ComponentFPosition           = "format string is at a wrong position"
	ErrF_PlaceholderArgumentMismatch = "mismatch between a placeholder and an argument: %v vs %v"
)

// NotificationTemplate is a template of a notification.
// A template describes contents of a notification and stores links to other
// sources of information.
type NotificationTemplate struct {
	name         Name
	components   []Component
	formatString *nm.FormatString
}

func NewNotificationTemplate(templateName string, args []any) (nt *NotificationTemplate, err error) {
	nt = &NotificationTemplate{
		name: cmb.Text(templateName),
	}

	nt.components, err = splitComponents(templateName, args)
	if err != nil {
		return nil, err
	}

	// If format string is present, set it as a separate field for fast access.
	for _, c := range nt.components {
		if c.Name == Component_F {
			nt.formatString, err = nm.NewFormatString((c.Value).(string))
			if err != nil {
				return nil, err
			}
		}
	}

	// Check the format string.
	if nt.formatString != nil {
		formatArgs := nt.components[1:]
		// Number and order of placeholders of the format string must be equal
		// to arguments for the format string.
		phs := nt.formatString.Placeholders()

		if len(phs) != len(formatArgs) {
			return nil, errors.New(Err_PlaceholdersCount)
		}

		for i, arg := range formatArgs {
			if arg.Name != phs[i].Type {
				return nil, fmt.Errorf(ErrF_PlaceholderArgumentMismatch, phs[i].Type, arg.Name)
			}
		}
	}

	return nt, nil
}

// splitComponents splits arguments into components and checks them.
func splitComponents(templateName string, args []any) (components []Component, err error) {
	// Fool check.
	if (len(templateName) == 0) || (len(args) == 0) {
		return nil, errors.New(Err_Stupid)
	}
	if len(templateName) != len(args) {
		return nil, errors.New(Err_ArgsCount)
	}

	// Check components' count.
	err = checkComponentsCount(templateName)
	if err != nil {
		return nil, err
	}

	components = make([]Component, 0, len(args))

	// Split the components.
	var cmp *Component
	var cmpName string
	for i, symbol := range ([]rune)(templateName) {
		cmpName = string(symbol)

		cmp, err = NewComponent(cmpName, args[i])
		if err != nil {
			return nil, err
		}

		components = append(components, *cmp)
	}

	// Check position of the format string if it is present.
	for i, c := range components {
		if c.Name == Component_F {
			if i != 0 {
				return nil, errors.New(Err_ComponentFPosition)
			}
		}
	}

	return components, nil
}
