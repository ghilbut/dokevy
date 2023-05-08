package storage

import (
	"path/filepath"
	"sort"
	"strings"

	// external
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"

	// project
	. "github.com/ultary/oss/xsr/internal/viewer/config"
)

type S3 struct {
	Client *s3.S3
	Bucket string
}

func NewS3(cfg *Config) Storage {
	var (
		accessKey = cfg.S3.AccessKey
		secretKey = cfg.S3.SecretKey
		token     = cfg.S3.Token
		region    = cfg.S3.Region
	)

	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, token),
		Region:      aws.String(region),
	})
	if err != nil {
		log.Fatal(err)
	}

	return &S3{
		Client: s3.New(sess),
		Bucket: cfg.Bucket,
	}
}

// Find 2 kind of s3 object key
//  1. exact matching media file - *.js, *.css, *.jpg, etc.
//  2. an index.html file that allows the SPA to serve path
func (s *S3) Find(host, path string) string {
	base := filepath.Join(host, path)
	if s.head(base) != "" {
		return base
	}

	index := filepath.Join(base, "index.html")
	if s.head(index) != "" {
		return index
	}

	keys := s.list(host)
	sort.Slice(keys, func(i, j int) bool {
		const sep = "/"
		ln := len(strings.Split(*keys[i], sep))
		rn := len(strings.Split(*keys[j], sep))
		return ln > rn
	})

	for _, key := range keys {
		if strings.HasPrefix(base, filepath.Dir(*key)) {
			return *key
		}
	}

	return ""
}

// head check there is an exact matching key or not
func (s *S3) head(key string) string {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.Client.HeadObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error(aerr.Error())
			}
		} else {
			log.Error(err.Error())
		}
		return ""
	}
	return key
}

// list looks for the index.html files belonging to the host
func (s *S3) list(host string) []*string {
	var (
		loop           = true
		next   *string = nil
		result         = make([]*string, 0)
	)

	for loop {
		const maxkeys = 1000
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.Bucket),
			ContinuationToken: next,
			MaxKeys:           aws.Int64(maxkeys),
			Prefix:            aws.String(host),
		}
		output, err := s.Client.ListObjectsV2(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case s3.ErrCodeNoSuchBucket:
					log.Error(s3.ErrCodeNoSuchBucket, aerr.Error())
				default:
					log.Error(aerr.Error())
				}
			} else {
				log.Error(err.Error())
			}
			return empty
		}

		for _, object := range output.Contents {
			key := object.Key
			if strings.HasSuffix(*key, "index.html") {
				result = append(result, key)
			}
		}

		loop = *output.IsTruncated
		next = output.NextContinuationToken
	}

	return result
}

// Serve response the s3 object file
func (s *S3) Serve(ctx *fasthttp.RequestCtx, key string) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	output, err := s.Client.GetObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Error(s3.ErrCodeNoSuchKey, aerr.Error())
			case s3.ErrCodeInvalidObjectState:
				log.Error(s3.ErrCodeInvalidObjectState, aerr.Error())
			default:
				log.Error(aerr.Error())
			}
		} else {
			log.Error(err.Error())
		}
		log.Fatal(err)
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType(*output.ContentType)
	if strings.HasSuffix(key, "index.html") {
		ctx.Response.Header.Set("Cache-Control", "max-age=0,must-revalidate,public")
	}
	ctx.SetBodyStream(output.Body, int(*output.ContentLength))
}
