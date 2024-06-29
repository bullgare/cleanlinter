module github.com/bullgare/cleanlinter/test/testdata

go 1.22.3

require github.com/bullgare/cleanlinter/test/testdata/src/project_correct v0.0.0-00010101000000-000000000000 // indirect
require github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect v0.0.0-00010101000000-000000000000 // indirect

replace github.com/bullgare/cleanlinter/test/testdata/src/project_correct => ./src/project_correct
replace github.com/bullgare/cleanlinter/test/testdata/src/project_incorrect => ./src/project_incorrect