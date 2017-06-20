# GOROOT
goroot="/usr/local/go"
if [ -d "$goroot" ]; then
    export GOROOT="$goroot"
fi
# PJROOT
PWDDIR=`pwd`
export PJROOT=$PWDDIR


# GOPATH
export GOPATH=$PJROOT
# PATH
export PATH=$PJROOT/bin:$GOROOT/bin:/bin:/sbin:/usr/sbin:/usr/bin:/usr/local/bin
