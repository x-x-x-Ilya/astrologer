
// Declarative //

pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                sh 'make build-app'
            }
        }
        stage('Lint') {
            steps {
                sh 'make lint'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Declarative Deploying log....'
            }
        }
    }
}