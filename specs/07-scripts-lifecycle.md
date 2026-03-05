# Lifecycle Scripts

## Description

k2 supports four script phases that allow executing shell commands at different points in the lifecycle of an apply. Scripts can be defined at two levels: in the **template** and in the **apply**.

## Phases

| Phase | Execution Time | Frequency |
|-------|----------------|-----------|
| `bootstrap` | Before any other action, during the very first apply | Once only |
| `pre` | Before file copying | On each apply |
| `post` | After file copying | On each apply |
| `nuke` | During destroy, before file deletion | On each destroy |

## Execution Order

During an **apply**:

```
1. [If first apply] template.scripts.bootstrap
2. [If first apply] apply.scripts.bootstrap
3. template.scripts.pre
4. apply.scripts.pre
5. ── Template file copying ──
6. template.scripts.post
7. apply.scripts.post
```

During a **destroy**:

```
1. template.scripts.post (if the template is resolved)
2. apply.scripts.nuke
3. ── Deletion of files listed in .gitignore ──
```

## First Apply Detection

An apply is considered "first" if:
- The target folder does not yet exist, **or**
- The `.gitignore` file does not exist in the target folder, **or**
- The `.gitignore` does not contain the line `!k2.apply`

## Script Syntax

Each script is a list of shell commands:

```yaml
scripts:
  bootstrap:
    - echo "first command"
    - echo "second command"
  pre:
    - mkdir -p ./data
  post:
    - echo "done {{ .name }}"
  nuke:
    - echo "cleanup {{ .name }}"
```

### Template Rendering in Scripts

Commands are rendered through the Go template engine before execution. Available variables depend on the context:

**For template scripts**: merged variables (`template.parameters` + `apply.vars`)

```yaml
# template k2.template.yaml
scripts:
  post:
    - echo "template fin of {{ .name }}"
```

**For apply scripts**: only apply variables (`apply.vars`)

```yaml
# k2.apply.yaml
scripts:
  bootstrap:
    - echo "boot {{ .name }}"
```

### Execution

Each script line is executed via `/bin/sh -c` in the apply folder (`apply.metadata.folder`):

```go
cmd := exec.Command("/bin/sh", "-c", renderedScript)
cmd.Dir = apply.K2.Metadata.Folder
```

The stdin, stdout, and stderr streams are connected to the current terminal.

### Empty Lines

Empty lines or lines containing only spaces are ignored.

## Examples (from samples)

### Scripts in a template (`kind1`)

```yaml
scripts:
  bootstrap:
    - echo "template boot of {{ .name }}"
  post:
    - echo "template fin of {{ .name }}"
```

### Scripts in an apply (`component1`)

```yaml
scripts:
  bootstrap:
    - echo "boot {{ .name }}"
  post:
    - echo "finjj"
```

### Nuke script in an apply (`component2`)

```yaml
scripts:
  nuke:
    - echo "Nuke script for component2"
```

### Bootstrap and post scripts in a Git apply (`componentFromGit1`)

```yaml
scripts:
  bootstrap:
    - echo "boot"
  post:
    - echo "fin"
```
