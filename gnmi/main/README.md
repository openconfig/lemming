## Running the Fake gNMI Server

```bash
go run main.go
```

```bash
gnmic -a localhost:1234 --insecure subscribe --mode stream --path openconfig:/system/state/current-datetime -u foo -p bar --target fakedut
```

For running on KNE (also experimental), see
https://github.com/wenovus/ondatra/tree/fake-prototype-0/fakebind
