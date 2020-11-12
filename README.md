
# grpcMiddlewareAuthentication  
This is an exemplary implementation of header-based authentication (JWT auth) using the [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware) authentication package. The primary objective here is to demonstrate the use of AuthFuncOverride so that the middleware makes an exception for the getToken methods that must not to be authenticated. 
Server and client implementation are with full code examples included, so please have a look.
Client calls the server to get a token and then calls with the issued token the greeter server.
  
### Run server

     go run cmd/server/server.go 

### Run client

      go run cmd/client/client.go 
  
## Contributing  
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.  
  
Please make sure to update tests as appropriate.  
  
## License  
[MIT](https://choosealicense.com/licenses/mit/)