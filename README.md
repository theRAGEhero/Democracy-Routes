# Democracy-Routes

Democracy Routes is software for discussion and decision-making.

## Setup

### Development Environment

To set up environment for development:

- Install [Go](https://go.dev/doc/install) programming language.
- Install [pre-commit](https://pre-commit.com/).

Then run:

```shell
make setup-dev-environment
```

### Development Infrastructure

To rise infrastructure for development:

- Install [rootless Docker](https://docs.docker.com/engine/security/rootless/).

Then use these commands to manage the infrastructure:

```shell
make dev-infra-start
```

```shell
make dev-infra-stop
```
