**netDep demo**

*URLs as Literals*
```sh
go build
./netDep -p ../demo/url-literal -s ../demo/url-literal -v -e "./env.yaml"
```

*Unresolved*
```sh
go build
./netDep -p ../demo/unresolved -s ../demo/unresolved -v 
```

*User-assisted*
```sh
go build
./netDep -p ../demo/annotated -s ../demo/annotated -v 
```
