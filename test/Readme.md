# It's all about integration testing

There are 2 test projects:
[project_correct](./testdata/src/project_correct) and [project_incorrect](./testdata/src/project_incorrect).

And 2 integration tests in [integration_test.go](./integration_test.go). They check if linter messages are triggered for correct lines (there are comments in projects' code like ` // want "Adapter layer is not allowed to reference Infrastructure layer"`).

\* Please don't bother about all the go modules. Unfortunately, it's a bare minimum for it to become working while not polluting the root module.