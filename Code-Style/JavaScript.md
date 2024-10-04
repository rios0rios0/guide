### General

#### Filenames
All filenames should use `snake_case`.

**Right:** `src/opensearch-dashboards/index_patterns/index_pattern.js`
**Wrong:** `src/opensearch-dashboards/IndexPatterns/IndexPattern.js`

#### Prettier and linting
We are gradually moving the OpenSearch Dashboards code base over to Prettier. All TypeScript code
and some JavaScript code (check `.eslintrc.js`) is using Prettier to format code. You
can run `node script/eslint --fix` to fix linting issues and apply Prettier formatting.
We recommend you to enable running ESLint via your IDE.

Whenever possible we are trying to use Prettier and linting overwritten developer guide rules.
Consider every linting rule and every Prettier rule to be also part of our developer guide
and disable them only in exceptional cases and ideally leave a comment why they are
disabled at that specific place.

### HTML
This part contains developer guide rules around general (framework agnostic) HTML usage.

#### Camel case `id` and `data-test-subj`
Use camel case for the values of attributes such as `id` and `data-test-subj` selectors.

```html
<button id="veryImportantButton" data-test-subj="clickMeButton">Click me</button>
```

The only exception is in cases where you're dynamically creating the value, and you need to use
hyphens as delimiters:

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

### TypeScript/JavaScript
The following developer guide rules apply for working with TypeScript/JavaScript files.

#### Prefer modern JavaScript/TypeScript syntax
You should prefer modern language features in a lot of cases, e.g.:

- Prefer arrow function over function expressions
- Prefer template strings over string concatenation
- Prefer the spread operator for copying arrays (`[...arr]`) over `arr.slice()`
- Use optional chaining (`?.`) and nullish Coalescing (`??`) over `lodash.get` (and similar utilities)

#### Avoid mutability and state
Wherever possible, do not rely on mutable state. This means you should not
reassign variables, modify object properties, or push values to arrays.
Instead, create new variables, and shallow copies of objects and arrays:

```js
// good
function addBar(foos, foo) {
  const newFoo = { ...foo, name: 'bar' };
  return [...foos, newFoo];
}

// bad
function addBar(foos, foo) {
  foo.name = 'bar';
  foos.push(foo);
}
```

#### Avoid `any` whenever possible
Since TypeScript 3.0 and the introduction of the
[`unknown` type](https://mariusschulz.com/blog/the-unknown-type-in-typescript) there are rarely any
reasons to use `any` as a type. Nearly all places of former `any` usage can be replaced by either a
generic or `unknown` (in cases the type is really not known).

You should always prefer using those mechanisms over using `any`, since they are stricter typed and
less likely to introduce bugs in the future due to insufficient types.

#### Use object destructuring
This helps avoid temporary references and helps prevent typo-related bugs.

```js
// best
function fullName({ first, last }) {
  return `${first} ${last}`;
}

// good
function fullName(user) {
  const { first, last } = user;
  return `${first} ${last}`;
}

// bad
function fullName(user) {
  const first = user.first;
  const last = user.last;
  return `${first} ${last}`;
}
```

#### Use array destructuring
Directly accessing array values via index should be avoided, but if it is
necessary, use array destructuring:

```js
const arr = [1, 2, 3];

// good
const [first, second] = arr;

// bad
const first = arr[0];
const second = arr[1];
```
