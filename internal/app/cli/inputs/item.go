package inputs

import (
	"github.com/manifoldco/promptui"

	"github.com/gsk148/gophkeeper/internal/app/validators"
)

func ItemID() (string, error) {
	ip := promptui.Prompt{Label: "Enter the item ID", Validate: validators.Min(1)}
	return ip.Run()
}

func ItemName() (string, error) {
	np := promptui.Prompt{Label: "Enter the item name", Validate: validators.ItemName}
	return np.Run()
}

func ItemNote() (string, error) {
	np := promptui.Prompt{Label: "Add a note (optional)", Validate: validators.Max(50)}
	return np.Run()
}

func ItemText() (string, error) {
	tp := promptui.Prompt{Label: "Enter the text", Validate: validators.Min(1)}
	return tp.Run()
}
