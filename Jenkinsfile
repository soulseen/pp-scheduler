pipeline {
    agent{
    kubernetes {
      label 's2i'
      yaml '''apiVersion: v1
kind: Pod
spec:
  containers:
  - name: ks-scheduler
    image: zhuxiaoyang/golang-1.12-ks-scheduler:latest
    command: [\'cat\']
    tty: true
    volumeMounts:
    - name: dockersock
      mountPath: /var/run/docker.sock
    - name: dockerbin
      mountPath: /usr/bin/docker
  volumes:
  - name: dockersock
    hostPath:
      path: /var/run/docker.sock
  - name: dockerbin
    hostPath:
      path: /usr/bin/docker
      '''
      defaultContainer 'ks-scheduler'
    }
}

    environment {
        DATA_PATH = '/tmp/test.db'
        GITHUB_CREDENTIAL_ID = 'github-id'
        KUBECONFIG_CREDENTIAL_ID = 'kubeconfig'
        DOCKERHUB_NAMESPACE = 'zhuxiaoyang'
        GITHUB_ACCOUNT = 'soulseen'
        KUBECONFIG = '/root/.kube/config'
    }

    stages {

        stage('set kubeconfig'){
         steps{
            sh 'mkdir -p ~/.kube'
            withCredentials([kubeconfigContent(credentialsId: "$KUBECONFIG_CREDENTIAL_ID", variable: 'KUBECONFIG_CONTENT')]) {
               sh 'echo "$KUBECONFIG_CONTENT" > ~/.kube/config'
            }
          }
        }

//        stage ('checkout scm') {
//            steps {
//                checkout(scm)
//            }
//        }

//        stage ('unit test') {
//           steps {
//                container ('go') {
//                    sh 'make test'
//                }
//            }
//        }

        stage ('e2e test') {
            steps {
                sh './hack/e2etest.sh'
            }
        }
    }
}