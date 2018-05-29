## Simple HTTP WebServer written on Golang
This is simple HTTP WebServer written on Golang to show GitLab continuous integration & deployment processes.  
GitLab will test and build main.go from any branch and deploy binary in test environment ('TestEnv' VMM, all brunches), and will deploy binary builded on deploy in production environment ('ProdEnv' VM, production brunch).   

## INSTALL GitLab Community Edition
See https://about.gitlab.com/installation/  
Update SMTP setting, see https://docs.gitlab.com/omnibus/settings/smtp.html  

## CREATE PROJECT FROM LOCAL FOLDER USING GIT
git config --global user.name "Yonas"  
git config --global user.email "blackover@yandex.ru"  
cd "C:\Users\Yonas\go\src\Simple-Golang-HTTP-Server"  
git init  
git remote add origin http://mygitlab.westus2.cloudapp.azure.com/root/Simple-Golang-HTTP-Server.git  
git git commit -a -m "Initial commit"  
git push -u origin master  
# to sync locally edited files from local to remote 
git add . 
git commit -m "v1"  
git push -u origin master  
# change branch to production  
git branch  
git checkout -b production  
git commit -a -m "production commit"  
git push origin production  

## CREATE .gitlab-ci.yml
```
# See https://docs.gitlab.com/ee/ci/yaml/
image: golang:latest

# we have to declare REPO_NAME variable to get working directory path
variables:
  REPO_NAME: Simple-Golang-HTTP-Server

# update working directory with source code
before_script:
  - ln -svf $CI_PROJECT_DIR /usr/lib/go-1.6/src
  - cd /usr/lib/go-1.6/src/$REPO_NAME

stages:
    - test
    - build
    - deploy

# job to test source code
format:
    stage: test
    tags:
      - test
    script:
# format source code
      - go fmt
# find suspicious construct
      - go vet
# run main_test.go
      - go test -race

# job to obtain binary
compile:
    stage: build
    tags:
      - test
    script:
# create binary
      - go build -race -ldflags "-extldflags '-static'" -o $CI_PROJECT_DIR/gohttp
    artifacts:
      paths:
# get (upload) created binary ("artifact")
        - gohttp

# enviroment for testing
testenv:
    stage: deploy
    tags:
      - test
    script:
# change service binary
      - sudo service gohttp stop
      - cp -fv /usr/lib/go-1.6/src/$REPO_NAME/gohttp /usr/local/bin/
      - sudo service gohttp start

# working enviroment
prodenv:
    stage: deploy
    tags:
      - prod
    script:
# change service binary
      - sudo service gohttp stop
      - cp -fv /usr/lib/go-1.6/src/$REPO_NAME/gohttp /usr/local/bin/
      - sudo service gohttp start
    only:
      - production
    dependencies:
# get binary from job 'compile'
      - compile
```

## INSTALL RUNNERS on VMS
See https://docs.gitlab.com/runner/install/linux-manually.html  
curl -L https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh | bash  
apt-get install gitlab-runner  
gitlab-runner register  
	...  
touch /usr/local/bin/gohttp  
chmod +x /usr/local/bin/gohttp  
chown gitlab-runner:gitlab-runner /usr/local/bin/gohttp  

## CREATE SERVICE
To run gohttp binary as a service:  
wget https://raw.github.com/frdmn/service-daemons/master/debian -O /etc/init.d/gohttp  
nano /etc/init.d/gohttp  
	...  
chmod +x /etc/init.d/gohttp  
chown gitlab-runner:gitlab-runner /etc/init.d/gohttp   
update-rc.d gohttp defaults  
 
cp -f /etc/sudoers.d/90-cloud-init-users /etc/sudoers.d/90-service-gohttp-users  
nano /etc/sudoers.d/90-service-gohttp-users  
	gitlab-runner ALL=(ALL) NOPASSWD: /usr/sbin/service gohttp *  
usermod -aG sudo gitlab-runner  
su gitlab-runner  
sudo service gohttp start  
