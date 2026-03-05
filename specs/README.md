# k2 — Technical Specifications

This folder contains the comprehensive technical documentation for all features of **k2**, a declarative YAML-driven template engine.

## Table of Contents

| Document | Description |
|----------|-------------|
| [01-architecture.md](01-architecture.md) | General architecture and core concepts |
| [02-inventory.md](02-inventory.md) | Inventory file (`k2.inventory.yaml`) |
| [03-templates.md](03-templates.md) | Template definition (`k2.template.yaml`) |
| [04-template-apply.md](04-template-apply.md) | Template application (`k2.apply.yaml`) |
| [05-template-sources.md](05-template-sources.md) | Template sources: inventory and git |
| [06-variables-and-rendering.md](06-variables-and-rendering.md) | Variables, template rendering and Sprig functions |
| [07-scripts-lifecycle.md](07-scripts-lifecycle.md) | Lifecycle scripts (bootstrap, pre, post, nuke) |
| [08-cli-commands.md](08-cli-commands.md) | CLI commands: plan, apply, destroy |
| [09-gitignore-and-files.md](09-gitignore-and-files.md) | `.gitignore` management and generated files |
| [10-deduplication.md](10-deduplication.md) | Task deduplication and hashing |
| [11-nested-templates.md](11-nested-templates.md) | Nested templates (template within a template) |
| [12-custom-files.md](12-custom-files.md) | Custom files co-located with applies |
