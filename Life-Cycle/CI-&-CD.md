## Context
This document is intend to describe the standards to follow creating a CI/CD pipeline.

## Standard
A perfect Continuous Integration and Continuous Deployment (CI/CD) pipeline would consist of the following steps:

1. **Source Control Management:** Code is stored in a source control management (SCM) system, such as Git, and committed to the repository by developers.
2. **Build:** The code is built, compiled and packaged into a deployable artifact by a build system such as Jenkins or TravisCI.
3. **Unit Testing:** Automated unit tests are run to ensure the code is functioning correctly and to catch any bugs or errors early on in the process.
4. **Integration Testing:** The code is integrated with other parts of the system and integration tests are run to ensure that everything is working together correctly.
5. **Static Analysis:** Code quality and security checks are performed using tools like SonarQube, ESLint, etc.
6. **Deployment:** The artifact is deployed to a staging environment for further testing and validation.
7. **Acceptance Testing:** Users and stakeholders test the system in the staging environment to ensure it meets their requirements.
8. **Release:** If the system passes all acceptance tests, it is released to production.
9. **Monitoring:** The system is monitored in production to ensure it is functioning correctly and to detect any issues that may arise.
10. **Continuous Feedback:** Feedback from users and stakeholders is used to improve the system and the pipeline.

Keep in mind that this is a high-level overview of a perfect CI/CD pipeline.

## References
