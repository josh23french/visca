pipeline {
  agent any
  tools {
    go 'go-1.16'
  }
  environment {
    GO111MODULE = 'on'
    GOPATH = '$(pwd)/go'
  }
  stages {
    stage('Build') {
      steps {
        sh 'go build'
      }
    }

    stage('Test') {
      steps {
        sh 'mkdir "$GOPATH"'
        sh 'go get github.com/axw/gocov/gocov'
        sh 'go get github.com/AlekSi/gocov-xml'
        sh '$GOPATH/bin/gocov test github.com/josh23french/visca | $GOPATH/bin/gocov-xml > coverage.xml'
        junit 'coverage.xml'
      }
    }

  }
}
