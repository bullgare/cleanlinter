# cleanlinter

Simplistic linter to check golang project internal imports according to Clean Architecture pattern.
In other workds, checks that layers are not referencing outwards.

To make it work, you can specify Clean Architecture layers (at least 2 out of 4):

**Domain** - usually contains domain layer business rules.

**UseCase** - usually contains services, reusable handlers, the app core logic goes here.

**Adapter** - usually contains adapters to other APIs.

**Infrastructure** - usually contains presentation logic (http-, grpc-handlers), DB integration, and so on.

You are doing it by specifying the layer's package path.

## Installation

```sh
go install github.com/bullgare/cleanlinter/cmd/cleanlinter@latest
```

## Usage

```sh
cleanlinter \
-cleanlinter_path_to_domain=github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/domain \
-cleanlinter_path_to_usecase=github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/usecase \
-cleanlinter_path_to_adapter=github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/adapter \
-cleanlinter_path_to_infra=github.com/bullgare/cleanlinter/test/testdata/src/project_correct/internal/infra \
./...
```

\* You don't have to specify all the paths, at least 2.

\** There is yet another flag, `-cleanlinter_verbose=true`, which will output more details, including found layers.

## Tests

### Unit

```sh
make test
```
### Integration

```sh
make integration-test
```