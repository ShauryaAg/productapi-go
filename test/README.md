## Testing

This is a little unconventional way of testing the API in Golang.

Under `test/main.go` a single **`test()`** function is defined that can be used to run all the test cases defined in the `<model>_cases.json` file.

Each test case contains the following fields:

```json
{
    "name": "",
    "endpoint": "",
    "method": "",
    "handler": "",
    "inputHeaders": {},
    "inputBody": {},
    "expected": {
        "status": 0,
        "response": {}
}
```

The `setUp()` function is used to set up the test cases, it creates a dummy user, product and review in the **`test`** DB and saves the user `jwtToken` and `id` for further test cases.
