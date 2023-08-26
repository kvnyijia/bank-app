package mail

import (
	"testing"

	"github.com/kvnyijia/bank-app/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGamil(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1>Hello world</h1>
		<p>hey hey hey</p>
	`
	to := []string{"????????@mail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
