# lambdapack

This project provides a Go based CLI which builds archives for Lambda functions, with the aim being to provide a simple way to migrate functions using `go1.x` lambda runtime to `provide.al2` by packaging existing binaries with a bootstrap script making it compatible with the newer runtime.

# Usage

Current lambda projects will build one or more binaries using the name of the parent folder, this could be the project name, or if your using the `./cmd` pattern the name of the folder containing the main package(s).

Often these binaries are output into a `bin/` folder, this CLI will package up each of these binaries along with a bootstrap script into a zip archive ready to deploy to AWS Lambda.

So given the following folder structure:

```
bin
├── api-lambda
└── scheduled-lambda
```

Running the following command:
```
lambdapack ./bin/* dist
```

This would build two zip archives:

```
dist
├── api-lambda.zip
└── scheduled-lambda.zip
```

Functions then reference these archives as follows.

```yaml
  LambdaFunction:
    Type: AWS::Serverless::Function
    Properties:
      Handler: bootstrap
      Runtime: provided.al2
      CodeUri: dist/example.zip
      Architectures:
        - arm64
```

# License

This project is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).