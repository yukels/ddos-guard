package awsclient

import (
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"

	"github.com/yukels/util/context"
)

type S3Client struct {
	s3session  *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3Client(ctx context.Context) (*S3Client, error) {
	if err := createSession(ctx); err != nil {
		return nil, err
	}

	return &S3Client{
		s3session:  s3.New(awsSession),
		uploader:   s3manager.NewUploader(awsSession),
		downloader: s3manager.NewDownloader(awsSession),
	}, nil
}

func (c *S3Client) GetS3API(ctx context.Context) s3iface.S3API {
	return c.s3session
}

func (c *S3Client) DownloadData(ctx context.Context, bucket, sourceFilePath string, writer io.WriterAt) error {
	if err := c.download(ctx, bucket, sourceFilePath, writer); err != nil {
		return errors.Wrapf(err, "Unable to download from [s3://%s/%s]", bucket, sourceFilePath)
	}
	return nil
}

func (c *S3Client) DownloadFile(ctx context.Context, bucket, sourceFilePath, destFilePath string) error {
	baseDir := filepath.Dir(destFilePath)
	err := os.MkdirAll(baseDir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Unable to create destination folder of file %s", destFilePath)
	}
	file, err := os.Create(destFilePath)
	if err != nil {
		return errors.Wrapf(err, "Unable to create file [%s]", destFilePath)
	}
	defer file.Close()

	if err := c.download(ctx, bucket, sourceFilePath, file); err != nil {
		return errors.Wrapf(err, "Unable to download from [s3://%s/%s]", bucket, sourceFilePath)
	}

	return nil
}

func (c *S3Client) download(ctx context.Context, bucket, sourceFilePath string, writer io.WriterAt) error {
	_, err := c.downloader.Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(sourceFilePath),
	})
	if err != nil {
		return err
	}

	return nil
}

type DownloadBatch struct {
	Name   string
	Bucket string
	Path   string
	Writer io.WriterAt
}

func (b *DownloadBatch) WriteAt(p []byte, off int64) (n int, err error) {
	return b.Writer.WriteAt(p, off)
}

func (c *S3Client) DownloadBatch(ctx context.Context, batch []*DownloadBatch) error {
	forDownload := make([]s3manager.BatchDownloadObject, 0, len(batch))
	for _, b := range batch {
		forDownload = append(forDownload, s3manager.BatchDownloadObject{
			Object: &s3.GetObjectInput{
				Bucket: aws.String(b.Bucket),
				Key:    aws.String(b.Path),
			},
			Writer: b.Writer,
		})
	}

	return c.downloader.DownloadWithIterator(ctx, &s3manager.DownloadObjectsIterator{Objects: forDownload})
}

func S3SplitPath(ctx context.Context, path string) (string, string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", "", errors.Wrapf(err, "Can't parse s3 path [%s]", path)
	}
	if strings.ToLower(u.Scheme) != "s3" {
		return "", "", errors.Wrapf(err, "Schema should be 's3' [%s]", path)
	}
	return u.Host, u.Path[1:], nil
}
