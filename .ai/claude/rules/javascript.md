---
paths:
  - "**/*.{js,jsx,ts,tsx}"
---

# JavaScript & TypeScript Conventions

> **TL;DR:** Use `snake_case` for file names and `camelCase` for HTML attribute values. Prefer modern syntax (arrow functions, template strings, optional chaining, nullish coalescing). Never use `any` -- use `unknown` instead. Favor immutability and destructuring.

## General

### File Names

All file names must use `snake_case`.

```
src/my_class.js    ✅ Correct
src/myClass.js       ❌ Wrong
```

### Formatting and Linting

Use **Prettier** for code formatting and **ESLint** for linting. Every Prettier and ESLint rule is considered part of this guide. Disable rules only in exceptional cases and always leave a comment explaining why.

## HTML

### Camel Case for `id` and `data-test-subj`

Use camelCase for the values of `id` and `data-test-subj` attributes:

```html
<button id="veryImportantButton" data-test-subj="clickMeButton">Click me</button>
```

The only exception is when dynamically generating values, where hyphens may be used as delimiters:

```jsx
buttons.map(btn => (
  <button
    id={`veryImportantButton-${btn.id}`}
    data-test-subj={`clickMeButton-${btn.id}`}
  >
    {btn.label}
  </button>
))
```

## TypeScript / JavaScript

### Prefer Modern Syntax

- Prefer **arrow functions** over function expressions.
- Prefer **template strings** over string concatenation.
- Prefer the **spread operator** (`[...arr]`) over `arr.slice()` for copying arrays.
- Use **optional chaining** (`?.`) and **nullish coalescing** (`??`) over `lodash.get` and similar utilities.

### Avoid Mutability

Do not reassign variables, modify object properties, or push values to arrays. Instead, create new variables and shallow copies:

```js
// ✅ Good
function addBar(foos, foo) {
  const newFoo = { ...foo, name: 'bar' };
  return [...foos, newFoo];
}

// ❌ Bad
function addBar(foos, foo) {
  foo.name = 'bar';
  foos.push(foo);
}
```

### Avoid `any`

Since TypeScript 3.0 introduced the [`unknown` type](https://mariusschulz.com/blog/the-unknown-type-in-typescript), there is rarely a valid reason to use `any`. Replace `any` with either a generic type parameter or `unknown`, combined with type narrowing.

### Use Object Destructuring

Destructuring reduces temporary references and prevents typo-related bugs:

```js
// ✅ Best
function fullName({ first, last }) {
  return `${first} ${last}`;
}

// ❌ Bad
function fullName(user) {
  const first = user.first;
  const last = user.last;
  return `${first} ${last}`;
}
```

### Use Array Destructuring

Avoid accessing array values by index. When direct access is necessary, use array destructuring:

```js
const arr = [1, 2, 3];

// ✅ Good
const [first, second] = arr;

// ❌ Bad
const first = arr[0];
const second = arr[1];
```

---

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
