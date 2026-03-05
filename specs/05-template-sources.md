# Template Sources

## Description

k2 supports two template resolution modes: **inventory** (local) and **git** (remote). The source is defined in the `template.source` field of the `k2.apply.yaml` file.

## `inventory` Source

### Description

Resolves a template from files already present in the local inventory, identified by its `id`.

### Configuration

```yaml
template:
  source: inventory
  params:
    id: <template identifier>
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | string | yes | The template's `metadata.id` in the inventory |

### Resolution Mechanism

1. The execution plan contains the list of all templates scanned in the inventory
2. Resolution looks for the template whose `metadata.id` matches the `id` parameter
3. The template is found via its SHA-256 hash computed from `source` + `params`
4. The template source folder is the one containing the corresponding `k2.template.yaml`

### Example

```yaml
# k2.apply.yaml
k2:
  metadata:
    id: my-service
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: k2.cli.sample.templates.kind1
    vars:
      name: my-component
```

## `git` Source

### Description

Clones a remote Git repository and extracts a template from a specific path within the repository.

### Configuration

```yaml
template:
  source: git
  params:
    repository: <repository url>
    branch: <branch>
    path: <path to k2.template.yaml in the repository>
```

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `repository` | string | yes | Git repository URL |
| `branch` | string | yes | Branch to clone |
| `path` | string | yes | Relative path to `k2.template.yaml` in the repository |

### Resolution Mechanism

1. A SHA-256 hash is computed from `source` + `params`
2. The repository is cloned into `.refs/<hash>/` (with `--single-branch` option)
3. If the `.refs/<hash>/.git` folder already exists, a `git pull` is performed instead of a clone
4. The template is read from `<path>` in the cloned repository
5. The template folder is the parent folder of the `k2.template.yaml` file in the clone

### Git Cache

Cloned repositories are stored in the `.refs/` folder at the inventory root:

```
my-project/
├── .refs/
│   └── a1b2c3d4.../           ← Template ref hash
│       ├── .git/
│       └── samples/templates/fromGit1/
│           ├── k2.template.yaml
│           └── README.md
```

During `destroy`, the `.refs/<hash>/` folders corresponding to Git templates are also deleted.

### Example

```yaml
# k2.apply.yaml
k2:
  metadata:
    id: k2.cli.sample.services.product2.componentFromGit1
    kind: template-apply
  body:
    template:
      source: git
      params:
        repository: https://github.com/tuxounet/k2.git
        branch: main
        path: samples/templates/fromGit1/k2.template.yaml
    vars:
      name: component2
      description: "Template of type kind1"
```

## Reference Hashing

The hash of a template reference is computed with SHA-256 on the concatenation of `source` and `params`:

```go
value := fmt.Sprintf("%s-%v", t.Source, t.Params)
sha256 := sha256.New()
sha256.Write([]byte(value))
hash := fmt.Sprintf("%x", sha256.Sum(nil))
```

This hash is used to:
- Uniquely identify a template reference
- Store Git clones in `.refs/<hash>/`
- Deduplication: two applies referencing the same template with the same params share the same resolution
