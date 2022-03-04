## Running the Fake gNMI Server

```bash
go run main.go
```

```bash
gnmic -a localhost:1234 --insecure subscribe --mode stream --path openconfig:/system/state/current-datetime -u foo -p bar --target fakedut
```
