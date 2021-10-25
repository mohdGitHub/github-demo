# Backend Test Project - Twitter Board

Your task is to implement a simple Twitter Board backend API. Detailed specifications for the test project are provided below. We estimate that you will not need more than a single weekend at relaxed coding speed to implement it.

## Project Description

The Twitter Board API will be used by your Users to perform the following tasks:

- User registration and login endpoints.
- Twitter Search endpoint (send a search query, and then it searches through API calls in twitter.com, and returns the first 50 results).
- Endpoint to save the previous results in a database.
- Endpoint that returns the results saved in the database, with pagination

## Technical details

Your backend should be able to serve all kinds of clients (which you do not have to implement) - using a RESTful API.

The following technical requirements are placed on your implementation:

### API

- Use Golang (v1.16+)
- HTTP responses should follow best practices
- API should communicate with their clients using JSON data structures
- Implement authentication that would be the best for the clients (using JWT)

### Data Storage

- All data should be stored in a relational database, use Postgres

### Users

- Registrations should be done with email and password
- You should implement the following functionality:
  - User Registration
  - User Login

### Twitter data

- You should implement the following functionality:
  - Twitter must have: description
  - Create a new Twitter
  - List all Twitters

### Test

- Your code should be tested using UnitTest
  - Models
  - Controllers 

### Bonus task (NOT mandatory)

- Endpoints that performs CRUD operation on the previous results
- Use GoLang concurrency in one of your functions

## Review process

There are a few technical restrictions, so we can see how you fare with the technologies and processes we use on a daily basis, but in general, the actual implementation is quite open-ended. The reason is we want to see how you think in terms of backend architecture, development processes and how you generally deal with the challenges you might face while implementing this app.

The following should help you determine where to put your focus, since these are the things we will be looking for while reviewing your project.

### 🔥 Code quality

Is your code well-structured? Do you keep your coding style consistent across your codebase?

### 🔥 Security

How do you store your customers' passwords? What about security of your customers' data? how you are securing the API endpoints

### 🔥 Testability

Is your code tested? How do you approach testing? Do you use TDD or are tests an afterthought for you?

### API structure and usability

How do you structure API endpoints? Do you follow REST principles? Do you make use of proper response codes and HTTP headers where it makes sense?

### Validations and error handling

How do you handle required fields, and errors that might appear due to invalid data,
How do you handle responses and Exceptions

### Development and deployment

How hard is it to run your project locally?

### Documentation

Is your API documented? Is your documentation auto generated from the code base? Does it cover all you endpoints?

### Version Control

Please commit often and tell a story of your process with your commit history.

## Project Delivery

- source code delivery options:
    * Create a Fork of this repository on Bitbucket and then create a pull request back to it
    * Create a private repository on bitbucket and invite hussein@abwaab.me

> That's it. Good luck, and we look forward to seeing your submission!
