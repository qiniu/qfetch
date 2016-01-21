export GOPATH=$GOPATH:/Users/jemy/QiniuCloud/Projects/qfetch
gox -os="darwin" -arch="amd64" -output="../bin/qfetch_darwin_amd64" 
gox -os="darwin" -arch="386" -output="../bin/qfetch_darwin_386"
gox -os="windows" -arch="amd64" -output="../bin/qfetch_windows_amd64"
gox -os="windows" -arch="386" -output="../bin/qfetch_windows_386"
gox -os="linux" -arch="amd64" -output="../bin/qfetch_linux_amd64"
gox -os="linux" -arch="386" -output="../bin/qfetch_linux_386"
gox -os="linux" -arch="arm" -output="../bin/qfetch_linux_arm"
