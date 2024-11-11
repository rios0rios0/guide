# Tips and Tricks
This page includes some of the tips and tricks while writing Python programs. They don’t quite deserve their own pages, so here they are, stuffed into this one page.

## Table of Contents
- [Use Descriptive Names](#use-descriptive-names)
- [Write Comments](#write-comments)

## Use Descriptive Names
Some developers like to over-abbreviate variable/class names:
```python
r = requests.get("https://example.com")
rj = r.json()
rc = r.status_code
```

These kinds of names are fine for PoC scripts and are commonly used in mathematical algorithms. However, a production project full of these names is not very readable to maintainers. Compare this code:
```python
q1 = queue.Queue()
q2 = queue.Queue()

while True:
    p = p1.get() * p2.get()
```

…with this:
```python
# these two queues are relational
customer_bill_queue = queue.Queue()
tax_rate_queue = queue.Queue()

while True:
    after_tax_rate = customer_bill_queue.get() * tax_rate_queue.get()
```

Using descriptive variable names can make your code much more readable.

## Write Comments
Other developers will have to reverse-engineer your code and guess what you were thinking when you were writing it to be able to maintain it. Instead of forcing other people to guess what you meant, write your thoughts down. You might even forget what you were thinking after a while, especially when the code is complex.
Instead of:
```python
for row in [r for r in results if r is not None]:
    parsed.append(dict(zip(field_name, row)))
```

…write this:
```python
# for every row in database query result that's not empty
for row in [r for r in results if r is not None]:

    # zip the field names and the row items and make the result into a dict
    # this is similar to using a MySQLCursorDict cursor
    # then append the dict into the parsed results list
    parsed.append(dict(zip(field_name, row)))
```

This is a bit of an extreme/made-up example, but you get the point. Although the code is still somewhat tangled, at least now you know what it is doing. You can move on if you don’t need to understand **how** it does its thing but only need to know **what** it does.
