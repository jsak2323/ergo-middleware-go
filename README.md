## INSTALL VM REQUIREMENTS

    export LANGUAGE=en_US.UTF-8
    export LC_ALL=en_US.UTF-8

    sudo apt-get update
    sudo apt-get -y upgrade

    sudo apt-get install software-properties-common
    sudo apt-get update && sudo apt-get install sqlite3
    sudo apt-get install jq
    sudo apt-get install zip
    sudo apt update && sudo apt -y install gcc g++

    // install stackdriver
    curl -sSO https://dl.google.com/cloudagents/add-monitoring-agent-repo.sh && sudo bash add-monitoring-agent-repo.sh && sudo apt-get update && sudo apt-cache madison stackdriver-agent

    sudo apt-get install -y 'stackdriver-agent=6.*'

    // version check
    dpkg-query --show --showformat '${Package} ${Version} ${Architecture} ${Status}\n' stackdriver-agent

## Setup MySql
    sudo apt install mysql-server
    sudo mysql_secure_installation
    sudo mysql -u root -p
    // run this query to require root password when connecting to db
    USE mysql; 
    UPDATE mysql.user SET plugin = 'mysql_native_password' WHERE user = 'root' AND host = 'localhost';
    ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'mynewpassword';
    FLUSH PRIVILEGES;
    exit;
    sudo service mysql restart

## NODE INSTALLATION
    mkdir ergo-node
    cd ergo-node
    wget https://github.com/ergoplatform/ergo/releases/download/v4.0.24/ergo-4.0.24.jar
    mkdir .ergo
    wget https://pastebin.com/raw/fez234Dy
    mv fez234Dy ergo.conf
    java -Xmx4G -jar ergo-4.0.24.jar --mainnet -c ergo.conf


## INSTALL GO

    cd /tmp
    wget https://dl.google.com/go/go1.14.7.linux-amd64.tar.gz

    sudo tar -xvf go1.14.7.linux-amd64.tar.gz
    sudo mv go /usr/local

    sudo nano ~/.bashrc
    // add to end of file
    export GOROOT=/usr/local/go
    export GOPATH=$HOME/go
    export PATH=$GOPATH/bin:$GOROOT/bin:$PATH

    source ~/.profile

    go version
  
## add nonce.sqlite to /var/db

## CLONE APP FROM GIT

    cd $GOPATH
    mkdir src && cd src
    mkdir github.com && cd github.com
    mkdir btcid && cd btcid

    git clone git@35.240.159.3:btcid/wallet/ergo-middleware-go.git
    cd ergo-middleware-go
    mkdir logs
    go mod init

## RUN TESTS

    go test ./...

## BUILD APP

    go build cmd/server/*.go

## RUN APP IN DEV MODE

    ./main

## RUN APP IN PRODUCTION MODE (adjust config.json before running in production mode)
    
    screen -S app
    PRODUCTION=true ./main


## RUN MIDDLEWARE RPC CLIENT

    go run client/*.go <command> 

    // in production
    PRODUCTION=true go run client/*.go <command> 

    // example: 
    go run client/*.go getbalance

    // in production
    PRODUCTION=true go run client/*.go getbalance



## API REFERENCE

    // getblockchaininfo
    curl -X 'GET' \
        'http://localhost:9052/info' \
        -H 'accept: application/json'
    


