## INIT GIT FROM LOCAL FOLDER
git config --global user.name "Yonas"  
git config --global user.email "blackover@yandex.ru"  
cd "C:\Users\Yonas\go\src\Simple-Golang-HTTP-Server"  
git init  
git remote add origin http://mygitlab.westus2.cloudapp.azure.com/root/Simple-Golang-HTTP-Server.git  
git add .  
git commit -m "Initial commit"  
git push -u origin master  
# to sync files from local to git after edit  
git add . 
git commit -m "v1"  
git push -u origin master  

## CREATE .gitlab-ci.yml
```
image: golang:latest
variables:
  REPO_NAME: Simple-Golang-HTTP-Server
before_script:
  - ln -svf $CI_PROJECT_DIR /usr/lib/go-1.6/src
  - cd /usr/lib/go-1.6/src/$REPO_NAME
stages:
    - test
    - build
    - deploy
format:
    stage: test
    tags:
      - test
    script:
      - go fmt
      - go vet
#       - go test -race $(go list ./... | grep -v /vendor/)
compile:
    stage: build
    tags:
      - test
    script:
      - go build -race -ldflags "-extldflags '-static'" -o $CI_PROJECT_DIR/gohttp
    artifacts:
      paths:
        - gohttp
run:
    stage: deploy
    tags:
      - test
    script:
      - sudo service gohttp stop
      - cp -f /usr/lib/go-1.6/src/$REPO_NAME/gohttp /usr/local/bin/
      - sudo service gohttp start
```

## INSTALL RUNNER
curl -L https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh | bash  
apt-get install gitlab-runner  
gitlab-runner register  
	...  
touch /usr/local/bin/gohttp  
chmod +x /usr/local/bin/gohttp  
chown gitlab-runner:gitlab-runner /usr/local/bin/gohttp  

## CREATE SERVICE
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
