# JavaScript & TypeScript Testing

> **TL;DR:** Use Jest for testing, Enzyme for React component rendering, and Faker for generating test data. All tests must follow the BDD pattern with `// given`, `// when`, `// then` comment blocks. Test business rules and user flows rather than framework internals. Use `mount` for integration tests (triggers `useEffect`) and `shallow` for lightweight unit tests.

## Overview

Front-end testing can be challenging. This document provides clear guidelines on **what** to test and **how** to test it, ensuring consistency across the codebase. It is a living document and should be updated as new modules and features introduce new testing requirements.

## Tools

| Tool                                         | Purpose                                 |
|----------------------------------------------|-----------------------------------------|
| [Jest](https://jestjs.io/)                   | JavaScript testing framework and runner |
| [Enzyme](https://enzymejs.github.io/enzyme/) | React component testing utility         |
| [Faker](https://github.com/marak/Faker.js/)  | Fake data generation library            |

## BDD Structure (Given / When / Then)

**Every test must use `// given`, `// when`, `// then` comment blocks** to clearly separate preconditions, actions, and assertions. This is mandatory across all test files.

```ts
it('should render Content when not loading', () => {
  // given
  const props = { loading: false };

  // when
  const component = shallow(<ExampleComponent {...props} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});
```

## Unit Tests

### What to Test

Focus on testing **your own** business rules and code behavior, not the framework's rendering internals:

**a. Props are passed correctly:**
```ts
it('should render an ExampleChildren with the right props', () => {
  // given
  const component = mount(<ExampleComponent prop={'prop_value'} />);

  // then
  expect(component.find(ExampleChildren).prop('prop')).toBe('prop_value');
});
```

**b. Component output matches expected state:**
```ts
it('should render the Content component', () => {
  // when
  const component = shallow(<ExampleComponent loading={false} />);

  // then
  expect(component.find(Content).exists()).toBeTruthy();
  expect(component.find(Loading).exists()).toBeFalsy();
});

it('should render Loading component', () => {
  // when
  const component = shallow(<ExampleComponent loading={true} />);

  // then
  expect(component.find(Content).exists()).toBeFalsy();
  expect(component.find(Loading).exists()).toBeTruthy();
});
```

**c. User interactions mutate state correctly:**
```ts
it('should present the Congrats component when button clicked', async () => {
  // given
  const component = mount(<ExampleComponent />);

  // when
  component.find(ExampleButton).simulate('click');
  component.update();

  // then
  expect(component.find(Congrats).exists()).toBeTruthy();
});
```

**d. Async actions execute correctly:**
```ts
import { flushComponent } from "tests/public/utils/flush_component.ts";

it('should render the Children with the loaded data', async () => {
  // given
  const exampleService = new MockExampleService();
  const exampleData = faker.random.alphaNumeric();
  exampleService.withSuccess(exampleData);
  const component = mount(
    <ExampleAsyncComponent />,
    { services: { exampleService } }
  );

  // when
  await flushComponent(component);

  // then
  expect(component.find(Children).exists()).toBeTruthy();
  expect(component.find(Children).prop('loadedData')).toBeTruthy();
});
```

**e. Logic returns expected output:**

Avoid duplicating business rules in tests. Hard-code inputs and expected outputs:
```ts
it('should return the sum multiplied by two', () => {
  // given
  const numbersCollection = new NumbersCollection([1, 2, 3]);

  // when
  const result = numbersCollection.sumAndDuplicate(); // (1 + 2 + 3) * 2

  // then
  expect(result).toBe(12);
});
```

**f. Router navigation (`useHistory`):**
```tsx
const mockHistoryPush = jest.fn();

jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useHistory: () => ({
    push: mockHistoryPush,
  }),
}));

it('should call push', async () => {
  // given
  const component = mount(<Something />);

  // when
  await flushComponent(component);

  // then
  expect(mockHistoryPush).toHaveBeenCalled();
});
```

### How to Test

| Utility         | Purpose                                                                                  |
|-----------------|------------------------------------------------------------------------------------------|
| `act`           | Required for any action that triggers React state updates                                |
| `flushPromises` | Forces Jest to resolve all pending promises (e.g., API requests)                         |
| `mount`         | Renders full DOM including `useEffect` hooks; suitable for integration tests             |
| `shallow`       | Renders a lightweight representation without `useEffect`; more performant for unit tests |

**Tip:** Use `console.log(component.debug())` to inspect the rendered output when debugging test failures.

## References

- [Jest Documentation](https://jestjs.io/docs/getting-started)
- [Enzyme Documentation](https://enzymejs.github.io/enzyme/)
- [React Testing Best Practices](https://kentcdodds.com/blog/common-mistakes-with-react-testing-library)
- [Given-When-Then -- Martin Fowler](https://martinfowler.com/bliki/GivenWhenThen.html)
