pipeline {
  agent any
  tools {
    go 'go-1.16'
  }
  environment {
    GO111MODULE = 'on'
    GOPATH = "${WORKSPACE}/go"
    PATH = "${PATH}:${WORKSPACE}/go/bin"
  }
  stages {
    stage('Build') {
      steps {
        sh 'go build'
      }
    }

    stage('Test') {
      steps {
        sh 'go get github.com/axw/gocov/gocov'
        sh 'go get github.com/AlekSi/gocov-xml'
        sh 'gocov test github.com/josh23french/visca | gocov-xml > coverage.xml'
        junit 'coverage.xml'
      }
    }

  }
}
