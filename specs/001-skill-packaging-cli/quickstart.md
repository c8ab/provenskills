# Quickstart: psk

## Prerequisites

- Go 1.23+ installed

## Build from source

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o psk ./cmd/psk/
```

## Package your first skill

1. Create a skill directory:

```bash
mkdir my-skill
cat > my-skill/SKILL.md << 'EOF'
---
name: my-skill
description: A sample skill for testing psk.
metadata:
  version: "1.0.0"
  author: "your-name"
---

# My Skill

Instructions for the agent go here.
EOF
```

2. Build the artifact:

```bash
psk build ./my-skill --maintainer "Your Name <you@example.com>"
```

Expected output:
```
Built skill: my-skill@1.0.0
  author:     your-name
  maintainer: Your Name <you@example.com>
  stored:     ~/.psk/store/my-skill/1.0.0/
```

3. List stored skills:

```bash
psk list
```

Expected output:
```
NAME      VERSION  AUTHOR     MAINTAINER
my-skill  1.0.0    your-name  Your Name <you@example.com>
```

4. Validate before building (optional):

```bash
psk validate ./my-skill
```

## Environment

Override the default store location:

```bash
export PSK_STORE=/path/to/custom/store
psk list
```

## Run tests

```bash
go test ./tests/... -v
```
