## Context
This document is intend to describe the standards to follow creating a CI/CD pipeline.

## Standard
There are many best practices for securing applications, but in our experience, those below need to be carefully considered:

1. **Input validation and sanitization:** This practice involves ensuring that all user input is properly validated and sanitized to prevent injection attacks such as SQL injection and cross-site scripting (XSS).
This can be accomplished through techniques such as data type checking, input length validation, and filtering or encoding special characters (OWASP Top 10 - A1).
2. **Authentication and Authorization:** Implementing robust authentication and authorization mechanisms can help ensure that only authorized users can access sensitive information.
This can include using multifactor authentication, implementing role-based access controls, and securing session management (OWASP Top 10 - A2).
3. **Encryption:** Encrypting sensitive data can help protect it from unauthorized access, both in transit and at rest.
This can include using secure protocols such as HTTPS for data in transit and using encryption for data stored on disk (OWASP Top 10 - A3).
4. **Access controls:** Implementing the least privilege and need-to-know access controls can help limit the exposure of sensitive data.
This can include implementing fine-grained access controls, using access control lists, and regularly reviewing and revoking access privileges (OWASP Top 10 - A4).
5. **Security testing:** Regularly performing security testing and vulnerability assessments can help identify and remediate potential threats.
This can include using both manual and automated testing methods, such as penetration testing and dynamic application security testing (DAST) (OWASP Top 10 - A5).
6. **Incident response:** Developing and implementing an incident response plan can help address security breaches effectively.
This can include identifying the scope of the incident, containing the damage, and implementing a plan for recovery (MITRE ATT&CK - T1571).
7. **Logging and monitoring:** Implementing logging and monitoring mechanisms can help detect and respond to potential security threats in real-time.
This can include using centralized logging and event management systems, as well as implementing intrusion detection and prevention systems (MITRE ATT&CK - T1028).
8. **Network security:** Implementing firewalls and other network security measures can help protect against network-based attacks.
This can include using network segmentation, using virtual private networks (VPNs), and implementing network access controls (MITRE ATT&CK - T1547).

## References:

* OWASP Top 10: https://owasp.org/top10/
* MITRE ATT&CK: https://attack.mitre.org/
* OWASP Proactive Controls: https://owasp.org/www-project-proactive-controls/
* Veracode Report: https://www.veracode.com/security/application-security-report
