package main
import(	
	"fmt"
  "mime/multipart"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
  "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/awserr"
)
var filepa string
var MyBucket string
func UploadImage(f multipart.File, objectid string, header *multipart.FileHeader) (string,error){
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(goDotEnvVariable("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			goDotEnvVariable("AWS_ACCESS_KEY_ID"),
			goDotEnvVariable("AWS_SECRET_ACCESS_KEY"),
			"", // a token will be created when the session it's used.
		  ),
	}))
    uploader := s3manager.NewUploader(sess)
    MyBucket = goDotEnvVariable("BUCKET_NAME")
    file := f
    filename := objectid
    //upload to the s3 bucket
    result, err := uploader.Upload(&s3manager.UploadInput{
    Bucket: aws.String(MyBucket),
    ACL:    aws.String("public-read"),
    Key:    aws.String(filename),
    Body:   file,
    })
    if err != nil {
		return "",err
    }
    filepa = "https://" + MyBucket + "." + "s3-" + goDotEnvVariable("AWS_REGION") + ".amazonaws.com/" + filename
    fmt.Println("file saved to S3: %v", filepa)
    return result.Location, nil
}

func DeleteImage(objectid string) error{
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String(goDotEnvVariable("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			goDotEnvVariable("AWS_ACCESS_KEY_ID"),
			goDotEnvVariable("AWS_SECRET_ACCESS_KEY"),
			"", // a token will be created when the session it's used.
		  ),
	}))
	svc := s3.New(sess)
	MyBucket = goDotEnvVariable("BUCKET_NAME")
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(MyBucket),
		Key:    aws.String(objectid),
	}
	result, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return err
	}
	fmt.Println(result)
	return nil
}
