# k2 (BETA)

K2 Build System CLI

## CLI

### from everywhere 

```bash
K2_VERSION=$(curl --silent "https://api.github.com/repos/tuxounet/k2/tags" | jq -r '.[0].name')
go install github.com/tuxounet/k2@${K2_VERSION}
$(go env GOPATH)/bin/k2 help
```


### install

```bash
go install github.com/tuxounet/k2@v0.9.0
$(go env GOPATH)/bin/k2 help
```

### uninstall

```bash
rm $(go env GOPATH)/bin/k2
```

## Usage

Allowed actions:

- plan (compute what action well be donc)

```bash
$(go env GOPATH)/bin/k2 plan
```

- apply (apply template directives)

```bash
$(go env GOPATH)/bin/k2 apply
```

- destroy (cleanup template files)

```bash
$(go env GOPATH)/bin/k2 destroy
```
