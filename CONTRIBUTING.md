# Contributing to Greye

Thank you for considering contributing to Greye! We welcome contributions from the community and are excited to see what you will bring to the project.


## Development Environment Setup

1. **Requirements**:
    - Go 1.21 or later

2. **Setup**:
   ```sh
   git clone https://github.com/your-username/greye.git
   cd greye
   go mod download
   ```

## How to Contribute

### Reporting Bugs

If you find a bug, please report it by opening an issue on our [GitHub repository](https://github.com/greye-monitoring/greye/issues). Include as much detail as possible to help us understand and reproduce the issue.

### Feature Requests

We welcome new feature requests! If you have an idea for a new feature, please open an issue on our [GitHub repository](https://github.com/greye-monitoring/greye/issues) and describe the feature in detail.

### Code Contributions

1. **Fork the Repository**: Fork the [repository](https://github.com/greye-monitoring/greye) to your own GitHub account.
2. **Clone the Repository**: Clone your forked repository to your local machine.
4. **Make Changes**: Make your changes to the codebase.
5. **Run Tests**: Ensure all tests pass before submitting your changes.
   ```sh
   go test ./...
   ```
6. **Commit Changes**: Commit your changes with a clear message.
   ```sh
   git commit -m "Add feature: your feature description"
   ```
7. **Push Changes**: Push your changes to your forked repository.
   ```sh
   git push origin feature/your-feature-name
   ```
8. **Submit a Pull Request**: Open a pull request from your branch to the main repository.

## Coding Standards

- Follow Go best practices and standard library conventions
- Add comments for public functions and complex logic
- Follow the project's architecture and design patterns

## Testing Guidelines

- Write unit tests for all new functionality
- Ensure all tests pass before submitting a pull request
- Add integration tests where appropriate

## Documentation

- Update the README.md if you change functionality
- Document new features with clear examples
- Update annotations documentation if you add new ones

## Review Process

All pull requests will be reviewed by the maintainers. We may suggest changes or improvements before merging. Please be patient during the review process.

## License

By contributing to Greye, you agree that your contributions will be licensed under the project's license.