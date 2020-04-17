S3_BUCKET=better
INPUT_FILE=template.yaml
OUTPUT_FILE=sam-template.yaml
STAGE_NAME=dev
STACK_NAME=better

#upload swagger file to S3 bucket
aws s3 cp open-api.yaml s3://$S3_BUCKET/open-api.yaml

#build
make build

#Package and upload to s3
aws cloudformation package --template-file $INPUT_FILE \
                          --output-template-file $OUTPUT_FILE \
                          --s3-bucket $S3_BUCKET

# deploy
aws cloudformation deploy --template-file $OUTPUT_FILE \
                          --stack-name $STACK_NAME
                          --capabilities CAPABILITY_IAM