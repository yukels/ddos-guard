package awsclient

import (
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

var awsSession *session.Session

func createSession(ctx context.Context) error {
	if awsSession != nil {
		return nil
	}

	var err error
	awsSession, err = session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return err
	}
	log.Log(ctx).Debug("aws session was created")
	return nil
}
