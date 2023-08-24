package mail

import (
	"testing"

	"github.com/Chien179/NMCBookstoreBE/src/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="https://github.com/Chien179">Minh Chien Simple Bank</a></p>
	`
	to := []string{"chiennm890@gmail.com"}
	// attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, nil)
	require.NoError(t, err)
}
