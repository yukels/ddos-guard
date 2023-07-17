package guard

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"

	"github.com/yukels/ddos-guard/config"
	awsclient "github.com/yukels/util/aws-client"
	"github.com/yukels/util/context"
	"github.com/yukels/util/http/request"
	"github.com/yukels/util/jwt"
	"github.com/yukels/util/log"
)

type UserService struct {
	config *config.UserServiceConfig
	client *awsclient.S3Client
}

type Users struct {
	WhiteListUsers   []string `json:"whiteListUsers,omitempty"`
	BlockedListUsers []string `json:"blockedListUsers,omitempty"`
}

func NewUserService(ctx context.Context, config *config.UserServiceConfig) (*UserService, error) {
	var client *awsclient.S3Client
	if config.S3Bucket != "" {
		var err error
		client, err = awsclient.NewS3Client(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "Can't initialize aws S3 client")
		}
	}
	s := &UserService{
		config: config,
		client: client,
	}

	return s, nil
}

func (s *UserService) Run(ctx context.Context) error {
	if s.config.S3Bucket == "" {
		return nil
	}

	s.getUsers(ctx)
	go s.usersLoop(ctx)
	return nil
}

func (s *UserService) GetBlockedUsers(ctx context.Context) []string {
	return s.config.BlockedListUsers
}

func (s *UserService) GetWhiteListUsers(ctx context.Context) []string {
	return s.config.WhiteListUsers
}

func (s *UserService) usersLoop(ctx context.Context) {
	log.Log(ctx).Info("User service thread is running...")
	waitPeriod := time.Duration(s.config.RefreshPeriod)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitPeriod):
			s.getUsers(ctx)
		}
	}
}

func (s *UserService) getUsers(ctx context.Context) {
	buffer := aws.NewWriteAtBuffer([]byte{})
	err := s.client.DownloadData(ctx, s.config.S3Bucket, s.config.S3Path, buffer)
	if err != nil {
		log.Log(ctx).WithError(err).Errorf("Can't refresh users from path: [s3://%s/%s]", s.config.S3Bucket, s.config.S3Path)
		return
	}

	users := Users{}
	data := buffer.Bytes()
	err = json.Unmarshal(data, &users)
	if err != nil {
		log.Log(ctx).WithError(err).Errorf("Can't parse users from path: [s3://%s/%s]", s.config.S3Bucket, s.config.S3Path)
		log.Log(ctx).WithError(err).Errorf("users: [%s]", string(data))
		return
	}

	s.config.WhiteListUsers = users.WhiteListUsers
	s.config.BlockedListUsers = users.BlockedListUsers
	log.Log(ctx).Debugf("White list users [%v]", s.config.WhiteListUsers)
	log.Log(ctx).Debugf("Blocked list users [%v]", s.config.BlockedListUsers)
}

func (s *UserService) UserFromRequest(ctx context.Context, req *http.Request) string {
	token, _ := request.TokenFromRequest(ctx, req)
	if len(token) == 0 {
		return ""
	}
	return jwt.UsernameFromUnverifiedToken(ctx, token)
}
