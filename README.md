# vk-status-stat


#### Installation Golang

    cd /tmp
    wget https://storage.googleapis.com/golang/go1.8.linux-amd64.tar.gz -nv
    sudo tar -xvf go1.8.linux-amd64.tar.gz
    sudo mv go /usr/local

#### Installation dependencies

    mkdir /var/work
    mkdir /var/work/go_libs

    export GOROOT=/usr/local/go
    export GOPATH=/var/work/go_libs:/var/work/vk-status-stat/
    export PATH=$PATH:$GOROOT/bin

    go get github.com/jinzhu/gorm
    go get github.com/deckarep/golang-set
    go get github.com/lib/pq