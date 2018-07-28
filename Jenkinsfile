node {
		def goImage
		def nodeImage

    stage('Clone Repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    stage('Build Images') {
        goImage = docker.build("shashwatdalal/go-lang-image","-f ./dockerfiles/Dockerfile.goLang .")
        nodeImage = docker.build("shashwatdalal/node-image","-f ./dockerfiles/Dockerfile.node .")

    }
    stage("Source Build / Tests") {
		 parallel (
          "GoLang": {
                   goImage.inside('-v $PWD:/go/src/go-ci-test') {
                     sh 'go build -o main'
                   }
            }
     )
    }

    stage('Push Image') {
        def prodImage = docker.build("shashwatdalal/prod-image","-f ./dockerfiles/Dockerfile.prod .")
        docker.withRegistry('https://registry.hub.docker.com', 'docker-hub-credentials') {
          prodImage.push("latest")
        }
    }

    stage('Deploy to Test') {
      node("integration-test") {
            sh 'sudo docker pull shashwatdalal/prod-image && \
                sudo docker stop main && \
                sudo docker run -d --rm --name main -p 80:8080 shashwatdalal/prod-image'
      }
    }

    stage('Regression Test') {
      node("katalon-chrome-tests") {
         sh 'cd ../../webapp-ui-tests && git pull && ./run_chrome'
      }
    }

    stage('Deploy to Production') {
			node("production") {
        sh 'sudo docker pull shashwatdalal/prod-image && \
            sudo docker stop main && \
            sudo docker run -d --rm --name main -p 80:8080 shashwatdalal/prod-image'
      }
    }
}