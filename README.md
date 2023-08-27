# lambdapack

This project provides a Go based CLI which builds archives for [AWS Lambda](https://aws.amazon.com/lambda/) functions, with the aim being to provide a simple way to migrate functions using `go1.x` lambda runtime to `provide.al2` by packaging existing binaries with a bootstrap script making it compatible with the newer runtime. 

For more details on the motivation and usage, please checkout my blog post on the deprecation of `go1.x` runtime [RIP AWS Go Lambda Runtime](https://www.wolfe.id.au/2023/08/09/rip-aws-go-lambda-runtime/?utm_source=github&utm_medium=lambdapack).

# Why?

I built this CLI as a way to do the absolute minimum required to migrate existing Go based lambda functions deployed via cloudformation templates. This is instead of using the SAM CLI as the basis for all my lambda projects. The SAM CLI is a great way to deploy your AWS serverless projects, but it is not a very good tool for building Go binaries for the following reasons:

1. This CLI has a lot of features related to building, packaging and deploying serverless applications but for a simple Go binary migration, it is overkill.
2. Given this project is focused on building for a large number of languages, both interpreted and compiled, it is not optimized for Go compilation.
3. Migrating existing projects built outside of SAM would require refactoring to use SAM specific project structures and deployment templates.
4. Debugging build failures in SAM is difficult due to the output being hidden in its own generic build process.

The goal of this CLI is to provide a way to apply minimal changes to existing Go binaries and package them up in a way that is compatible with the newer `provide.al2` runtime without needing to refactor existing projects or build processes. To do this I add a small bootstrap script that runs the existing binary, whatever its name is, this removes the need to rename everything to bootstrap.

This simple bootstrap script solution was suggested by @aidansteele, initially I discounted it, but alas after trialing SAM for a project I thought better of my choice and wrote this CLI.
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