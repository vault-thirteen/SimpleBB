package s

import (
	"errors"

	c "github.com/vault-thirteen/SimpleBB/pkg/common"
)

// MessageSettings are settings of e-mail messages.
type MessageSettings struct {
	SubjectTemplateForRegVCode string `json:"subjectTemplateForRegVCode"`
	SubjectTemplateForReg      string `json:"subjectTemplateForReg"`
	BodyTemplateForRegVCode    string `json:"bodyTemplateForRegVCode"`
	BodyTemplateForReg         string `json:"bodyTemplateForReg"`
	BodyTemplateForLogIn       string `json:"bodyTemplateForLogIn"`
	BodyTemplateForPwdChange   string `json:"bodyTemplateForPwdChange"`
	BodyTemplateForEmailChange string `json:"bodyTemplateForEmailChange"`
}

func (s MessageSettings) Check() (err error) {
	if (len(s.SubjectTemplateForRegVCode) == 0) ||
		(len(s.SubjectTemplateForReg) == 0) ||
		(len(s.BodyTemplateForRegVCode) == 0) ||
		(len(s.BodyTemplateForReg) == 0) ||
		(len(s.BodyTemplateForLogIn) == 0) ||
		(len(s.BodyTemplateForPwdChange) == 0) ||
		(len(s.BodyTemplateForEmailChange) == 0) {
		return errors.New(c.MsgMessageSettingError)
	}

	return nil
}
