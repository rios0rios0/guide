# JavaScript & TypeScript Conventions

> **TL;DR:** Use `snake_case` for file names and `camelCase` for HTML attribute values. Prefer modern syntax (arrow functions, template strings, optional chaining, nullish coalescing). Never use `any` -- use `unknown` instead. Favor immutability and destructuring.

## General

### File Names

All file names must use `snake_case`.

```
src/opensearch-dashboards/index_patterns/index_pattern.js    ✅ Correct
src/opensearch-dashboards/IndexPatterns/IndexPattern.js       ❌ Wrong
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

## References

- [Prettier](https://prettier.io/)
- [ESLint](https://eslint.org/)
- [TypeScript `unknown` Type](https://mariusschulz.com/blog/the-unknown-type-in-typescript)
- [MDN - Destructuring Assignment](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Destructuring_assignment)
