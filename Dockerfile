FROM scratch

WORKDIR $GOPATH/src/github.com/aeon27/myblog
COPY . $GOPATH/src/github.com/aeon27/myblog

EXPOSE 8000
CMD ["./myblog"]