## GraphQL Compatibility Acceptance Tests

[![Join the chat at https://gitter.im/graphql-cats/graphql-cats](https://badges.gitter.im/graphql-cats/graphql-cats.svg)](https://gitter.im/graphql-cats/graphql-cats?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

The aim of this project is to provide a set of compatibility acceptance tests for libraries 
that implement [GraphQL specification](https://github.com/facebook/graphql). Since GraphQL 
server implementations are written in different programming languages, this project defines test cases in a 
language independent format ([YAML](http://yaml.org)).

Using this test suite has following advantages:

* The users and the author of a GraphQL library will get a better feeling whether library is compliant to the GraphQL 
  specification and if not, which parts of the library are not compliant.
* It reduces amount of work a library implementor needs to do
* This makes it much easier for existing implementations to keep up with the specification changes 

### Contribution & Open Questions

At the moment this project finds itself in an infant phase. Any contributions are very welcome! Please join us in 
a [gitter chat](https://gitter.im/graphql-cats/graphql-cats). If you would like to 
join this effort and be part of the project, please let us know.

There is still a number of open questions that we need to figure out together:

* [File format and structure](https://github.com/graphql-cats/graphql-cats/issues/3)
* [Versioning & Distribution](https://github.com/graphql-cats/graphql-cats/issues/4)
* [Licence & Contributor Code of Conduct](https://github.com/graphql-cats/graphql-cats/issues/2)

### Repository format

`scenarios` folder contains a set of `*.yaml` and `*.graphql` files. Every subfolder represents a particular part of GraphQL specification. 
`*.yaml` files contain scenarios with a list of test cases. `*.graphql` files contain IDL schema definitions that are used in some of the 
scenarios (they are explicitly referenced). 

### Scenario File Format

Every scenario is a [YAML](http://yaml.org) file with following structure: 

* `scenario` - _String_ - the name of this scenario
* `background` - _Object_ (optional) - common definitions used by all of the tests
  * `schema` - _String_ (optional) - inline GraphQL IDL schema definition
  * `schema-file` - _String_ (optional) - IDL schema definition file path relative to the scenario file 
  * `test-data` - _Object_ (optional) - test data used for query execution and directives 
* `tests` - _Array of Objects_ - list of tests
  * `name` - _String_ - a name of the test
  * `given` - _Object_ - input information for the test
    * `query` - _String_ - the GraphQL query to execute an action against
    * `schema` - _String_ (optional) - inline GraphQL IDL schema definition
    * `schema-file` - _String_ (optional) - IDL schema definition file path relative to the scenario file
    * `test-data` - _Object_ (optional) - test data used for query execution and directives     
  * `when` - _Object_ - action that should be performed in the test. See the **Actions** section for a list of available actions.
  * `then` - _Object_ | _Arrays of Objects_ - assertions that verify result of an action. See the **Assertions** section for a list of available actions.

Definitions in the `given` part of a test may override definitions defined in the `background` section.
    
#### Actions

* **Query validation**
  * `validate` - _Array of Strings_ - the list of validation rule names to validate a query against. This action will only validate query without executing it. 
* **Query execution**
  * `execute` - _Object_ - executes a query
    * `validate-query` - _Boolean_ (optional) - disables query validation during the execution
    * `test-value` - _String_ (optional) - the name of a field defined in the `test-data`. This value would be passed as a root value to an executor.  
    
#### Assertions

* **Validation passes**
  * `passes` - _Any_ - verifies that validation was successful. Only applicable in conjunction with query validation action  
* **Error count**
  * `error-count` - _Number_ - number of the errors in execution/validation results  
* **Error contains match**
  * `error` - _String_ - execution/validation results contain provided error message (provided error message may contain only part of the actual message)  
  * `loc` - _Array of Objects_ | _Array of Arrays of Numbers_ (optional) - a list of error locations
    * `line` - _Number_ 
    * `column` - _Number_ 
* **Error regex match**
  * `error-regex` - _String_ - execution/validation results contain provided error message (uses provided regular expressions to match an error message)  
  * `loc` - _Array of Objects_ | _Array of Arrays of Numbers_ (optional) - a list of error locations
    * `line` - _Number_ 
    * `column` - _Number_ 