# Task Deduplication

## Description

The execution plan automatically deduplicates identical tasks to avoid resolving the same template multiple times or executing redundant actions.

## Mechanism

### Uniqueness Calculation

Each task is identified by a key composed of its type and parameters:

```go
actionKey := fmt.Sprintf("%s-%v", action.Type, action.Params)
```

### Deduplication Process

1. The plan iterates through all tasks
2. For each task, it computes its unique key
3. If the key does not yet exist, the task is kept
4. If the key already exists, the task is ignored (duplicate)
5. The relative order of retained tasks is preserved

### Deduplication Timing

Deduplication is applied at the end of the `Plan()` phase, after all tasks have been added.

## Use Cases

### Shared Templates

If multiple applies reference the same template (same source and same parameters), the template resolution is performed only once:

```yaml
# component1 → template kind1
# component2 → template kind1
# component2bis → template kind1
```

The three applies generate the same `local-resolve` task with the same hash, but it appears only once in the plan.

### Reference Hashing

The SHA-256 hash of the template reference ensures uniqueness:

```go
value := fmt.Sprintf("%s-%v", t.Source, t.Params)
hash := sha256(value)
```

Two references are considered identical if and only if they have the same source AND the same parameters.

## Example

With the samples, the raw plan would contain:

```
local-resolve  hash=AAA  (for kind2)
local-resolve  hash=BBB  (for kind1, requested by component2)
local-resolve  hash=BBB  (for kind1, requested by component2bis)  ← duplicate
local-resolve  hash=BBB  (for kind1, requested by kind2/sub)      ← duplicate
git-resolve    hash=CCC  (for fromGit1)
apply          component1
apply          component2
apply          component2bis
apply          componentFromGit1
apply          withCustomFiles
```

After deduplication:

```
local-resolve  hash=AAA  (for kind2)
local-resolve  hash=BBB  (for kind1)
git-resolve    hash=CCC  (for fromGit1)
apply          component1
apply          component2
apply          component2bis
apply          componentFromGit1
apply          withCustomFiles
```
