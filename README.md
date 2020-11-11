# grpcMiddlewareAuthentication
This is an exemplary implementation of using the [go-grpc-middleware] (https://github.com/grpc-ecosystem/go-grpc-middleware) authentication package for header-based authentication. The primary objective here is to demonstrate the use of AuthFuncOverride so that the middleware makes an exception for the getToken methods and does not need to be authenticated. This would make sense if there is a "getToken" service on the same server as other authenticated services.


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
