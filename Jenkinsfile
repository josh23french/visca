pipeline {
  agent any
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
        sh '$GOPATH/bin/gocov test github.com/josh23french/visca | $GOPATH/bin/gocov-xml > coverage.xml'
        junit 'coverage.xml'
      }
    }

  }
}