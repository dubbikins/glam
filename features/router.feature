Feature: Router
    In order to efficient serve http request 
    As a developer
    I need to be able to create a declarative and composible router that implements the http.Handler interface
    
    Scenario: Add a GET handler to a root Router
        Given there is a root Router
        And the router has a "GET" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "GET",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        
    Scenario: Add a PUT handler to a root Router
        Given there is a root Router
        And the router has a "PUT" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "PUT",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a POST handler to a root Router
        Given there is a root Router
        And the router has a "POST" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "POST",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a Delete handler to a root Router
        Given there is a root Router
        And the router has a "DELETE" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "DELETE",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a Patch handler to a root Router
        Given there is a root Router
        And the router has a "PATCH" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "PATCH",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a Connect handler to a root Router
        Given there is a root Router
        And the router has a "CONNECT" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "CONNECT",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a Head handler to a root Router
        Given there is a root Router
        And the router has a "HEAD" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "HEAD",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a Options handler to a root Router
        Given there is a root Router
        And the router has a "OPTIONS" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "header": {
                },
                "method": "OPTIONS",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Add a Trace handler to a root Router
        Given there is a root Router
        And the router has a "TRACE" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
            """
            {
                "method": "TRACE",
                "path": "/test",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """

    Scenario: Mount a subrouter with GET Handler
        Given there is a root Router
        And a subrouter with a "GET" handler mounted at "/sub" for path "/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
             """
            {
                "path": "/sub/test",
                "method": "GET",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Mount a subrouter with prefix PUT Handler
        Given there is a root Router
        And a subrouter with a "PUT" handler mounted at "/sub" for path "/{test}" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
             """
            {
                "path": "/sub/test",
                "method": "PUT",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: Mount a subrouter with regex POST Handler
        Given there is a root Router
        And a subrouter with a "POST" handler mounted at "/sub" for path "/(test:^abc123$)" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
             """
            {
                "path": "/sub/abc123",
                "method": "POST",
                "body": "success"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
    Scenario: A strict, param, and regex have overlapping matches and the correct precedence is applied
        Given there is a root Router
        And the router has a "GET" handler for path "/sub/test" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success-strict"
            }
            """
        And the router has a "GET" handler for path "/sub/{param}" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 300,
                "body": "success-param"
            }
            """
        And the router has a "GET" handler for path "/sub/(regex:^[a-zA-Z]*[\d]*$)" that responds with:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 400,
                "body": "success-regex"
            }
            """
        
        When a request is made:
             """
            {
                "path": "/sub/abc123",
                "method": "GET"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 400,
                "body": "success-regex"
            }
            """
        When a request is made:
             """
            {
                "path": "/sub/123abc",
                "method": "GET"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 300,
                "body": "success-param"
            }
            """
        When a request is made:
             """
            {
                "path": "/sub/test",
                "method": "GET"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "content-type": "text/plain"
                },
                "statusCode": 200,
                "body": "success-strict"
            }
            """
    Scenario: a middleware is added that updates the response header
        Given there is a root Router
        And the router has middleware that adds the following to the response header:
            """
            {
                "Authorization": "Basic abc123"
            }
            """
        And the router has a "GET" handler for path "/test" that responds with:
            """
            {
                "header": {
                    "Content-Type": "text/plain"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
        When a request is made:
                """
            {
                "path": "/test",
                "method": "GET"
            }
            """
        Then the Router responds
        And the response should match:
            """
            {
                "header": {
                    "Content-Type": "text/plain",
                    "Authorization": "Basic abc123"
                },
                "statusCode": 200,
                "body": "success"
            }
            """
            

