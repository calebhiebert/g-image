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
        echo 'Build Complete'
      }
    }

    stage('Push Image') {
      steps {
        withCredentials([usernamePassword(credentialsId: 'docker-hub-login', passwordVariable: 'PASSWORD', usernameVariable: 'USERNAME')]) {
          sh "docker login -u ${env.USERNAME} -p ${env.PASSWORD}"
          sh "docker push panchem/gfile"
        }
      }
    }
  }
}