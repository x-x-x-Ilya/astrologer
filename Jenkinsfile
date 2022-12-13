
// Declarative //

pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Declarative Building..'
            }
        }
        stage('Test') {
            steps {
                echo 'Declarative Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Declarative Deploying....'
            }
        }
    }
}

// Script //

node {
    stage('Build') {
        echo 'Script Building....'
    }
    stage('Test') {
        echo 'Script Testing....'
    }
    stage('Deploy') {
        echo 'Script Deploying....'
    }
}