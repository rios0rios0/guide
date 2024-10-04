[UNDER CONSTRUCTION]

- [Getting started guide](#getting-started-guide)
  - [Key technologies](#key-technologies)
  - [Fork and clone OpenSearch Dashboards](#fork-and-clone-opensearch-dashboards)
  - [Bootstrap OpenSearch Dashboards](#bootstrap-opensearch-dashboards)
  - [Run OSD Plugin](#run-osd-plugin)
- [Code guidelines](#code-guidelines)
  - [General](#general)
  - [TypeScript/JavaScript](#typescriptjavascript)

## Getting started guide

This guide is for any developer who wants a running local development environment where you can make, see, and test changes. It's opinionated to get you running as quickly and easily as possible, but it's not the only way to set up a development environment.

### Key technologies

OpenSearch Dashboards is primarily a Node.js web application built using React. To effectively contribute you should be familiar with HTML ([usage guide](#html)), SASS styling ([usage guide](#sass-files)), TypeScript and JavaScript ([usage guide](#typescriptjavascript)), and React ([usage guide](#react)).

```bash
$ git clone git@github.com:opensearch-project/OpenSearch-Dashboards.git
$ git clone https://xpto.com/xpto/osd-plugin OpenSearch-Dashboards/plugins/osd-plugin
```

### Bootstrap OpenSearch Dashboards

If you haven't already, change directories to your cloned repository directory:

```bash
$ cd OpenSearch-Dashboards
```

The `yarn osd bootstrap` command will install the project's dependencies and build all internal packages and plugins. Bootstrapping is necessary any time you need to update packages, plugins, or dependencies, and it's recommended to run it anytime you sync with the latest upstream changes.

```bash
$ yarn osd bootstrap
```

Note: If you experience a network timeout while bootstrapping:

```
| There appears to be trouble with your network connection. Retrying...
```

You can run command with —network-timeout flag:

```
$ yarn osd bootstrap —network-timeout 1000000
```

Or use the timeout by configuring it in the [`.yarnrc`](https://github.com/opensearch-project/OpenSearch-Dashboards/blob/main/.yarnrc). For example:

```
network-timeout 1000000
```

If you've previously bootstrapped the project and need to start fresh, first run:

```bash
$ yarn osd clean
```

### Run OSD Plugin

Start the OpenSearch Dashboards development server:

```bash
$ yarn start
```

When the server is up and ready (the console messages will look something like this),

```
[info][listening] Server running at http://localhost:5603/pgt
[info][server][OpenSearchDashboards][http] http server running at http://localhost:5603/pgt
```

click on the link displayed in your terminal to
access it.

Note - it may take a couple of minutes to generate all the necessary bundles. If the Dashboards link is not yet accessible, wait for a log message like

```
[success][@osd/optimizer] 28 bundles compiled successfully after 145.9 sec, watching for changes
```

Note: If you run a docker image, an error may occur:

```
Error: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]
```

This error is because there is not enough memory so more memory must be allowed to be used:

```
$ sudo sysctl -w vm.max_map_count=262144
```

For windows:

```
$ wsl -d docker-desktop
$ sysctl -w vm.max_map_count=262144
```

## Code guidelines

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
