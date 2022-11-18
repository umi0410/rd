```shell
export GOPATH=${HOME}/go
go install github.com/spf13/cobra-cli@latest
```

```shell
cobra-cli init --author "Jinsu Park dev.umijs@gmail.com" --license apache --viper
cobra-cli add create
cobra-cli add delete
cobra-cli add list
cobra-cli add reload
cobra-cli add run 
```
