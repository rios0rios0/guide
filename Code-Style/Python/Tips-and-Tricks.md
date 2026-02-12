# Tips and Tricks

> **TL;DR:** Use descriptive variable and class names to maximize readability. Write comments that explain **why** the code exists, not just **what** it does.

## Overview

This page collects practical tips for writing clean, maintainable Python code. These guidelines complement the other Python standards documents.

## Use Descriptive Names

Avoid over-abbreviated variable and class names. While short names are acceptable for prototypes or mathematical algorithms, production code must be readable by any team member.

**Poorly named:**
```python
r = requests.get("https://example.com")
rj = r.json()
rc = r.status_code
```

**Well named:**
```python
response = requests.get("https://example.com")
response_body = response.json()
status_code = response.status_code
```

A more complex comparison:

```python
# ❌ Cryptic
q1 = queue.Queue()
q2 = queue.Queue()

while True:
    p = q1.get() * q2.get()
```

```python
# ✅ Self-documenting
customer_bill_queue = queue.Queue()
tax_rate_queue = queue.Queue()

while True:
    after_tax_amount = customer_bill_queue.get() * tax_rate_queue.get()
```

## Write Meaningful Comments

Other developers (and your future self) will need to understand your code. Write comments that capture the **intent** behind non-obvious logic:

```python
# ❌ No context
for row in [r for r in results if r is not None]:
    parsed.append(dict(zip(field_name, row)))
```

```python
# ✅ Clear intent
# Filter out empty rows from the database query results
for row in [r for r in results if r is not None]:
    # Combine field names with row values into a dictionary
    # (similar to using a MySQLCursorDict cursor)
    parsed.append(dict(zip(field_name, row)))
```

The code may still be dense, but now a reader can understand **what** it does without reverse-engineering the logic, and can move on if they only need to know the intent rather than the implementation details.

## References

- [PEP 8 -- Naming Conventions](https://peps.python.org/pep-0008/#naming-conventions)
- [PEP 257 -- Docstring Conventions](https://peps.python.org/pep-0257/)
