
build:
	mkdir -p functions
	go get github.com/aws/aws-lambda-go/events
	go get github.com/aws/aws-lambda-go/lambda
	go get github.com/aws/aws-lambda-go/lambdacontext
	go get github.com/awslabs/aws-lambda-go-api-proxy/core
	go get github.com/fauna/faunadb-go/v4/faunadb
	go get github.com/mailjet/mailjet-apiv3-go
	go get github.com/dgrijalva/jwt-go
	GOBIN=${PWD}/functions go install ./...

