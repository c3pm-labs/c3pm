package input

import "github.com/AlecAivazis/survey/v2"

var loginSurvey = []*survey.Question{
	{
		Name:     "login",
		Prompt:   &survey.Input{Message: "Email or username"},
		Validate: survey.Required,
	},
	{
		Name:     "password",
		Prompt:   &survey.Password{Message: "Password"},
		Validate: survey.Required,
	},
}

type LoginPayload struct {
	Login    string
	Password string
}

func Login() (*LoginPayload, error) {
	payload := &LoginPayload{}
	err := survey.Ask(loginSurvey, payload, SurveyOptions...)
	return payload, err
}
