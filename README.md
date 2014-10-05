dap
===

Drag file to process it with shell command.

![demo.png](https://raw.github.com/doloopwhile/dap/master/demo.png)

`dap` command show a window and simple print path of dropped file.

## Install
```sh
go get github.com/doloopwhile/dap
```

## Example
Copy files to backup by drag and drop
```sh
dap | xargs cp {} $HOME/backup
```

Upload files to S3 by drag and drop
```sh
dap | xargs aws s3 cp {} s3://your-s3-bucket/dir/
```
