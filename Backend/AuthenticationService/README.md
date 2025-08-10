# Authentication service 
This Microservice handles User Authentication , token and session Management 

### Core Responsibilities
- Login.
    - Compares incoming credentials with stored credentials.
    - On successful login, an access and refresh token pair is created and returned.
    - Access token returned as response body.
    - Refresh token sent as a cookie.
- Logout.
    - Removes the Refresh token from the client.
    - Adds refresh token to the block list.
- Registration.
    - Validates request credentials.
    - Hashes the Password before storing it in the database.
- Access and refresh token:
    - Adds request refresh token to the Redis blocklist.
    - Generates new access and refresh tokens and returns them to the user.


### Tech stack used 
- Language: GoLang
- Database: Postgres
- Blocklist Store: Redis
- Password Hashing: Scrypt
- Access and Refresh tokens: JWT   
    - Asymmetric signing using RSA
    - Symmetrical signing using HMAC

### Folder Structure
``` bash
├── App
├── Certificates
│   ├── asymmetricKeys
│   └── symmetricKey
├── Config
├── Controllers
├── docs
├── dtos
├── Helpers
│   ├── Password
│   │   └── Password_test
│   ├── Response
│   └── Token
│       └── Token_test
├── Middleware
├── Models
│   └── Models_test
├── Routes
├── Services
└── Utils
```

- #### App
    - Application entry point and core application setup.
- #### Certificates
    - Stores asymmetric PEM files, used for the asymmetric JWT signing algorithm.
- #### Config 
    - Configuration structs and logic for  PostgreSQL, Redis and other external services.
- #### Controllers 
    - Controllers to handle individual route requests and responses
- #### Docs
    - Contains documentation for the current api version.
- #### Dtos
    - Data transfer objects, primarily for post operations and data validation.
- #### Helpers
    - #### Password  
        - Password Hash and comparison functions.
    - #### Response 
        - Contains functions to format HTTP responses.
    - #### Token
        - Contains logic for the creation, validations, parsing and blocking of refresh and access tokens.
- #### Middleware 
    - Middleware used for routes.
- #### Models 
    - Contains Models and query logic.
- #### Routes
    - Api endpoint declarations and route grouping.
- #### Services 
    - Contains Services used by their controller counterpart.
- #### Utils
    - Generic utility Functions.

For full api endpoint analysis, look at the ```api_documentation.yml``` file located in the ```/docs``` directory of each api version. 

## Future additions.
- Message Broker implementation for inter-service communication (Kafka or Redis Streams).
- Account deletion.
- Account Information Updating.
- The password hash output string will store the configurations used to create the hash.