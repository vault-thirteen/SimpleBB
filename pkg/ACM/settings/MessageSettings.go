package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// MessageSettings are settings of e-mail messages.
type MessageSettings struct {
	SubjectTemplate            string `json:"subjectTemplate"`
	BodyTemplateForReg         string `json:"bodyTemplateForReg"`
	BodyTemplateForLogIn       string `json:"bodyTemplateForLogIn"`
	BodyTemplateForPwdChange   string `json:"bodyTemplateForPwdChange"`
	BodyTemplateForEmailChange string `json:"bodyTemplateForEmailChange"`
}

func (s MessageSettings) Check() (err error) {
	if (len(s.SubjectTemplate) == 0) ||
		(len(s.BodyTemplateForReg) == 0) ||
		(len(s.BodyTemplateForLogIn) == 0) ||
		(len(s.BodyTemplateForPwdChange) == 0) ||
		(len(s.BodyTemplateForEmailChange) == 0) {
		return errors.New(c.MsgMessageSettingError)
	}

	return nil
}
