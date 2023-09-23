# 0.1
* Read manual config files
* Simple 200 tests for verification (response code)
* Return codes for 0 if ok 1 if test failed, 2 every place else
* Test 401 and 403 if in config files
* Test body for similar content
* Test headers for content if needed

# 0.2
* Checking url for documentation and checking format
* Build config for testing
* Change request data to see how api act on unsure data
* Create variables if set to string 255 - test string 256 254 etc.

# 0.3
* Structure of response to check
* Compare two url with response (for same structure)

# 0.4
* Multiple jwt tokens for multiple users

# 0.5
* Developer mode website on port 9821
* On developer mode can make request on web and get results there, but saved into DB
* Database for storing results - no sql like redis
* Developer history for checking
* Developer show url that failed and places where the problem is

# 1.0
* test all code
* add test to test if test are testing ok
* howto write manual configs, more documentation