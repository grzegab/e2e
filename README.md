# E2E test app
This app is for testing api.

## How does it works

## Ignore 500 errors
App test all different status codes and body. By default, 2xx and 4xx are tested. 
5xx errors are ignored because they indicate server error not user error. 
If user enter wrong data he should receive 4xx error.

### Grouping of request
Because of order for requests and ID needed for some endpoints groups were introduced. 
Each group has verbs (methods) to go with. Order of methods to execute: POST, GET, PATCH, PUT, DELETE and others if are any.

// url =>
//      ID: nil
//		verbs =>
//			POST: ...
//			PATCH: ...
//			PUT: ...
//			DELETE: ...
//			GET: ...
//		templateID =>
//			ID: 123
//			verbs =>
//				POST:...
//				DELETE:...
//			someUrl =>
//				ID: ...
//				POST: ..


## Usage
docker run --rm e2e e2e-tests {-l}

Option -l is for starting web server with manual testing option

Initial commit:
Greg "Grzegab" Gabryel