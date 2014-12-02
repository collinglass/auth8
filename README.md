# auth8

Token Authentication in Go.

### Step 1

Try to ```GET http://localhost:1337/api/users```

You will get Response:

```
{
    "Error": "Unauthorized",
    "Data": null,
    "Status": 401
}
```

### Step 2

Create an account with ```POST http://localhost:1337/api/auth/signup```

With Request Body:

```
{
    "email": "example@example.com",
    "password": "example"
}
```

You will get a similar Response: 

```
{
    "Error": "",
    "Data": {
        "Password": "example",
        "OtherData": "User Data",
        "Email": "example@example.com",
        "Id": "BpLnfgDsc2"
    },
    "Status": 200
}
```

### Step 3

Login with ```POST http://localhost:1337/api/auth/login```

With Request Body:

```
{
    "email": "example@example.com",
    "password": "example"
}
```

You will get a similar Response: 

```
{
    "Error": "",
    "Data": {
        "Token": "WD8F2qNfHK"
    },
    "Status": 200
}
```

### Step 4

Create a bearer token in the Authorization header: ```Authorization: Bearer WD8F2qNfHK```.

Try Step 1 again.

You will get Response:

```
{
    "Error": "",
    "Data": {
        "Users": [
            {
                "Password": "",
                "OtherData": null,
                "Email": "",
                "Id": ""
            },
            {
                "Password": "",
                "OtherData": null,
                "Email": "",
                "Id": ""
            },
            {
                "Password": "",
                "OtherData": null,
                "Email": "",
                "Id": ""
            },
            {
                "Password": "",
                "OtherData": null,
                "Email": "",
                "Id": ""
            },
            {
                "Password": "",
                "OtherData": null,
                "Email": "",
                "Id": ""
            },
            {
                "Password": "example",
                "OtherData": "User Data",
                "Email": "example@example.com",
                "Id": "BpLnfgDsc2"
            }
        ]
    },
    "Status": 200
}
```