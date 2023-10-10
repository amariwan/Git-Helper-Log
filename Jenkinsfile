node('BuildDebian9') {
    stage('GIT Checkout') { // for display purposes
        deleteDir()
        // Get some code from a GitHub repository
        git branch: 'develop', url: 'git@gitmaster.hq.aland-mariwan.de:tools/githelperlog.git'
        sh 'githelper clone'
        stash includes: '**/*', name: 'gitCheckoutArea', useDefaultExcludes: false
        // Get the Maven tool.
        // ** NOTE: This 'M3' Maven tool must be configured
        // **       in the global configuration.
        goHome = tool name: 'Go', type: 'go'
    }
}

parallel Debian9: {
node('BuildDebian9') {
    stage('Preparation') { // for display purposes
        deleteDir()
        unstash name: 'gitCheckoutArea'
        goHome = tool name: 'Go', type: 'go'
    }

    stage('Build') {
        // Run the maven build
        withEnv(["GOROOT=$goHome", "PATH+GO=${goHome}/bin"]) {
            sh '''
                go get -d ./...
                go build -o githelperlog
            '''
        }
    }
    stage('smokTest') {
        // Run the maven build
        withEnv(["GOROOT=$goHome", "PATH+GO=${goHome}/bin"]) {
            sh '''
                go get -u github.com/stretchr/testify/assert
                go test -run Unit

            '''
        }
    }
    stage('integration tester') {
        // Run the maven build
        withEnv(["GOROOT=$goHome", "PATH+GO=${goHome}/bin"]) {
            sh '''
                go get -u github.com/stretchr/testify/assert
                go test -run Integration

            '''
        }
    }

    stage('Results') {
        archiveArtifacts 'githelperlog'
    }
}
}, VS2017: {
node('vs2017') {
    stage('Preparation') { // for display purposes
        deleteDir()
        unstash name: 'gitCheckoutArea'
        goHome = tool name: 'Go', type: 'go'
    }

    stage('Build') {
        // Run the maven build
        withEnv(["GOROOT=$goHome", "PATH+GO=${goHome}/bin"]) {
            bat('''
                go get -d ./...
                go build -o githelperlog.exe
            ''')
        }
    }
     stage('smokTest') {
        // Run the maven build
        withEnv(["GOROOT=$goHome", "PATH+GO=${goHome}/bin"]) {

            bat( '''
                go get -u github.com\\stretchr\\testify\\assert
                go get -u "github.com\\traherom\\memstream"
                go get -u "github.com\\apsdehal\\go-logger"
                go test -run Unit

            ''')

        }
    }
    stage('integration tester') {
        // Run the maven build
        withEnv(["GOROOT=$goHome", "PATH+GO=${goHome}/bin"]) {

            bat( '''
                go get -u github.com\\stretchr\\testify\\assert
                go test -run Integration

            ''')

        }
    }
    //input 'press Proceed to continue'
    stage('Results') {
        archiveArtifacts 'githelperlog.exe'
    }
}}
