package inputs

import (
	"github.com/manifoldco/promptui"

	"github.com/gsk148/gophkeeper/internal/app/validators"
)

func FilePath() (string, error) {
	pp := promptui.Prompt{Label: "Enter the file path", Validate: validators.Min(5)}
	return pp.Run()
}
