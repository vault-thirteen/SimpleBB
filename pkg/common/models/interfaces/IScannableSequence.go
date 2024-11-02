package interfaces

type IScannableSequence interface {
	IScannable
	Next() bool // = HasNextValue()
}
