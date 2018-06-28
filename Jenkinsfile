pipeline {
  agent any

  stages {
    stage('Build Image') {
      agent {
        dockerfile {
          filename 'Dockerfile'
          additionalBuildArgs '-t panchem/gfile'
        }
      }

      steps {
        echo 'Build Dockerfile'
      }
    }
  }
}