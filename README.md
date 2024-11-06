## **GitHub Actions Workflow**

This repository uses a **GitHub Actions workflow** to automate testing and reporting. The workflow runs tests both in parallel (using matrix strategy) and separately for fixed input tests. Here’s a breakdown of what the workflow does:

### **Workflow Overview**
1. **Triggers:**
   - **Scheduled Runs**: The workflow runs every Monday at 12:00 PM.
   - **Push Events**: Triggered when code is pushed to the `main` branch.
   - **Pull Requests**: Runs on pull requests targeting the `main` branch.
   - **Manual Trigger**: Can also be triggered manually through GitHub.

2. **Parallel Fuzz Testing**:
   - The workflow uses **matrix strategy** to divide the fuzz tests into two parts:
     - One job tests inputs from range 0–5.
     - Another job tests inputs from range 5–10.
   - Each job uploads its results as artifacts.

3. **Fixed Input Test**:
   - Runs a separate test with a hardcoded input to verify expected output.
   - The results of this test are also uploaded as an artifact.

4. **Report Aggregation**:
   - After all the tests complete, the workflow downloads the individual reports.
   - It **aggregates all reports** into a single file.
   - The aggregated report is uploaded as a final artifact.

### **Key Features**
- **Automated Setup**:
  The workflow **automatically sets up the Go environment**, ensuring consistency and reliability across test runs.
  
- **Flexible Triggers**:  
  The workflow can be triggered through **scheduled events, push or pull requests**, or manually via the GitHub interface.
  
- **Parallel Execution with Go Routines and Channels**:  
  In this project, **Go routines** are used to run tests concurrently, improving efficiency by allowing multiple tests to execute simultaneously. **Channels** facilitate communication between the routines by collecting reports.

- **Artifact Uploads**:  
  All individual test reports are uploaded as artifacts, ensuring that results are available for inspection separately from the aggregated one.

- **Intermediate Reports Access**:  
  Developers can **access intermediate reports** directly from the **GitHub Actions dashboard** to inspect specific ranges or input test results. This is useful for troubleshooting individual tests or checking results early once one parallel job has finished.

- **Aggregated Report**:  
  Once all tests are complete, the individual reports are combined into a **single summary report** for easy review and visibility.

## **How to Extend**
1. **Add New Tests**:  
   You can create new test cases by defining functions using Go's `testing` package. For example, add new tests in the `tests` directory that test different inputs or edge cases.

2. **Modify the Fuzzer**:  
   You can modify the fuzzer's input generation logic to test different scenarios or improve coverage. Edit the `fuzzer.go` file to introduce new types of inputs or test more complex cases.

3. **Update Workflow**:  
   Modify the `.github/workflows/fuzz-testing.yml` file to change test parameters, adjust the test ranges, or customize how tests are triggered. You can also extend the workflow to run tests on multiple Go versions or environments.

## **How to Run**

### **Run the Program Using Docker**

#### **Why Use Docker?**
Using Docker ensures that you **don't need to install dependencies like OpenSSL** or TPM-related libraries on your local machine. Docker containers provide a **consistent and isolated environment**, making it easier to run the program without worrying about dependency conflicts or differences across platforms (like macOS, Windows, or Linux). This approach avoids the complexity of manually managing dependencies, such as:
- OpenSSL libraries for cryptographic operations.
- Build tools like `gcc` for compiling certain dependencies.
- TPM libraries for interacting with the TPM simulator or hardware.

> **Note:** If you prefer running the program locally, you would need to install Go and TPM libraries manually. However, using Docker is the recommended and easiest way to avoid dependency issues.

With Docker, everything the program needs is **packaged into a container**, ensuring anyone can run it with a single command, regardless of the underlying OS or environment.

1. **Build the Docker Image:**
   ```bash
   docker build -t my-little-fuzz-tester-go .
   ```

2. **Run the Docker Container:**
   ```bash
   docker run --rm my-little-fuzz-tester-go
   ```

3. **Expected Output:**
   You should see something like:
   ```bash
   Random Bytes: eaf5afde6fccabcd1234567890abcdef
   ```

### **How to run the test files (Unit tests)**

**Run the tests:**
   ```bash
   go test
   ```

If all tests pass, the output will look like this:
   ```bash
   PASS
   ok  	my-little-fuzz-tester-go	0.123s   
   ```

## **Contributing**
Contributions are welcome! If you find a bug or have an idea for an improvement, feel free to open an issue or submit a pull request.

