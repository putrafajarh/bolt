# `pkg` Directory

The `pkg` directory is a conventional location in Go projects to store reusable code that can be imported by other projects or internal packages. It is a best practice to use this directory for code that is not specific to the application but provides functionality that can be shared across multiple projects.

## Purpose

The `pkg` directory is used to:
- Organize reusable libraries and helper functions.
- Separate application-specific code from generic, reusable code.
- Promote modularity and maintainability.

## Best Practices

1. **Keep It Reusable**  
    Code in the `pkg` directory should be designed to be reusable and not tightly coupled to the application logic.

2. **Use Clear and Descriptive Names**  
    Subdirectories should have clear and descriptive names that reflect their purpose. For example:
    ```
    pkg/
    ├── logger/       # Logging utilities
    ├── config/       # Configuration management
    ├── middleware/   # HTTP middleware
    └── utils/        # General utility functions
    ```

3. **Avoid Overloading**  
    Do not put all code into the `pkg` directory. Application-specific code should reside in the `internal` or `cmd` directories.

4. **Document Each Package**  
    Provide a `README.md` or inline comments for each subdirectory to explain its purpose and usage.

5. **Follow Go Naming Conventions**  
    Use lowercase, short, and meaningful names for packages. Avoid underscores or mixed case.

6. **Test Thoroughly**  
    Write unit tests for all reusable code to ensure reliability and maintainability.

7. **Minimize Dependencies**  
    Keep dependencies minimal to make the code easier to reuse in other projects.

## Example Structure

```plaintext
pkg/
├── auth/          # Authentication utilities
│   ├── jwt.go     # JWT-related functions
│   └── oauth.go   # OAuth-related functions
├── db/            # Database utilities
│   ├── connection.go
│   └── migrations/
├── logger/        # Logging utilities
│   └── logger.go
└── utils/         # General-purpose utilities
     ├── strings.go
     └── time.go
```

## References

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Effective Go](https://go.dev/doc/effective_go)

By following these best practices, the `pkg` directory can serve as a robust foundation for reusable and maintainable Go code.