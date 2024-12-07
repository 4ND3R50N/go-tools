# go-tools
This library contains useful helper functions to do things with lesser code. Examples: Filter(), Map(), ToPointer(), Any().

## Helper functions
- Mapper
    - Transformer functions (e.g.: Mapping Structs)
    - Easier to read
    - Reducing boilerplate code
- Converter
    - e.g.: Transform values to pointer
- Filter
    - Iteration functions that the `slices` package is not providing yet
    - Easier to read
    - Reducing boilerplate code
- comparer
  - Iterate through a list and check if anything matches
  - Iterate through a list and check if all matches
- dbutils
  - Contains NewPGXLocks to do advisory locking
  - Contains Transaction that wrapps pgx to do transactions + locking