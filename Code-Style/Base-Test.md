## Context
The layers mentioned below are defined and exemplified inside this page [here](../Life-Cycle/Architecture/Backend-Design.md).

## File Structure
The default file structure is described in the [backend design section](../Life-Cycle/Architecture/Backend-Design.md).

## Commands
For each listener in the command, it is necessary to create a test.
ALWAYS describe in the method, what is the listener being tested. Or use some string to describe.

Use the provided description field to write as follows: `"should call <LISTENER> when ..."` and the rest will be the complement of what the method does.

## Controllers
For each status code, it is necessary to create a test to verify if the HTTP code is correct and if it is really returning the data that is being returned is correct.
ALWAYS describe in the method, what is the answer expected. Or use some string to describe.

Use the provided description field to write as follows: `"should respond <HTTP_STATUS_CODE> (HTTP_STATUS) when ..."` and the rest will be the complement of what the method does.

## Services
For each public method it's necessary AT LEAST a pair of tests (success and failure).
ALWAYS describe in the method, what is the method tested and if it's a "successful" return or an "error". Or use some string to describe.

Use the provided description field to write as follows: `"should ... when ...")` and the rest will be the complement of what the method does.

## Repositories
Exactly the same as in [Services](#services).

## Doubles
Martin Fowler, a software developer and author, has written about the different types of test doubles, which are objects that can be used in place of real objects during testing.
* A **faker** is a test double that generates fake data for use in tests.
* A **mock** is a test double that is configured to expect certain method calls and can be used to verify that the expected calls were made during the test.
* A **stub** is a test double that returns a pre-configured response when a certain method is called. It can be used to isolate the test subject from external dependencies.
* A **dummy** is a test double that is passed around but never actually used. It is typically used to fill an argument that is required by the test subject but is not actually used.
* **In-memory** is a term used to refer to the fact that the test double is stored in memory during the test and not persisted to disk or a database.

Martin Fowler's book "Refactoring: Improving the Design of Existing Code" in which he defined and discussed about test doubles such as Fakers, Mocks, Stubs, Dummies and In-memory.

So, a **double** is any object that is the actual object for testing, used to simulate an external dependency.
From that definition and the previous Martin Fowler's definition, we just use some of those things, like:
* **Stub:** just a canned answer, don't do memory logic here.
* **Dummy:** ready-made answers if possible, such as empty lists and null answers.
* **In-memory:** do in-memory logic here without using other modules.
* **Faker:** we use as an external library to generate smart faked data for us.
* **Mock:** we want to avoid them as possible. But if we can't, we use as an external library to mimic and verify method calls.

## Builders
The Builder Design Pattern is a creational design pattern that allows for the construction of complex objects step by step through a builder object.
It separates the construction of a complex object from its representation and allows for the same construction process to create different representations.
The builder pattern is a design pattern that allows for the step-by-step construction of complex objects using a builder object.

We usually use this pattern to create the classes for the automated testing.

## References

* https://martinfowler.com/bliki/TestDouble.html
* https://martinfowler.com/articles/mocksArentStubs.html
