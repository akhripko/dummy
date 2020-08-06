# dummy
Golang dummy http service

#gRPC
brew install protobuf

go get google.golang.org/grpc

go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway

#kafka
https://github.com/confluentinc/confluent-kafka-go

brew install librdkafka pkg-config

#git:mod
git config --global url."git@github.com:".insteadOf "https://github.com/"

git config --global url."git@gitlab.com:".insteadOf "https://gitlab.com/"

git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"

#graphql
https://gqlgen.com/getting-started/

$ mkdir gqlgen-todos

$ cd gqlgen-todos

$ go mod init github.com/[username]/gqlgen-todos

$ go get github.com/99designs/gqlgen