package models

type IScannableSequence interface {
	IScannable
	Next() bool // = HasNextValue()
}
