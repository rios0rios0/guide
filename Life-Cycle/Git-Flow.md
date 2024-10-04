## Context
This document describes the ways to work with git-flow, versioning, and how to work with branches.

## Flow
Git flow is a feature branch model, with these characteristics:

1. The `main` is always in a deployable state, that is, it can go to production.
2. All the changes will be added to the `main` through feature branches and pull requests.
3. All feature and fix branches will be merged to the `main` through `merge/pull requests`.
4. Any synchronization or conflict resolution should be made using the `git rebase` command. You can visit [this tutorial](http://atlassian.com/git/tutorials/rewriting-history/git-rebase) to have a better understanding of `git rebase` and the reason why using this is important.

![](Assets/feature-branches.svg)

The flow to create some modification of the code must follow the steps listed below.

1. Create **a new feature branch** with the pattern:
   ```
   feat/TASK-ID
   ```
2. While developing the feature, will be committed to the branch created for this purpose.
3. Create builds for this feature in the development environment, if needed.
4. When the feature development is complete, first of all, will need to **rebase with the main** to update your branch and resolve conflicts, if exists.
5. **Create a Pull Request in GitLab for the default branch** (usually `main`), adding your QA team to review.
6. Analyses the code, if all is correct, your PR is approved.
7. After your pull request is approved, merge your feature branch with the default branch.
8. The automated pipeline will generate the build of the branch `main`, and it'll delivery this build to the QA environment.
9. When the QA deployed version is approved by the acceptance team, the developr will follow the instructions in the `CHANGELOG.md` file, inside the repository to generate a production release.
10. Once the instructions are completly followed, the automated pipeline will deliver the release to the production environment.

## Naming Convention

### Branches
The names of the branches created must follow the following naming pattern:

If **there is** a task created in the project's issue tracker, the card id must be used: `type/TASK-ID`.

Examples:
```
feat/RD-000
fix/SG-990
```

If **there is NOT** a card in the issue tracker, use the change scope: `type/scope`.

Examples:
```
feat/add-logs
fix/input-mask
```

The **type** must be related to what that branch seeks to solve, for example, if a new feature will be developed, if a correction will be implemented, etc.

Some possible types are as follows:

* feat (a new feature will be implemented);
* fix (a fix will be implemented to resolve an existing issue);
* refactor (some refactor will be performed on part of the code that already works);
* chore (will implement some code and/or infrastructure improvement);
* test (some new test scenario will be developed for the application);
* docs (some documentation change will be made).

### Commit Messages
Below have how you need to write a better commit message in the pattern.
1. Always use the `task-id`. If there is no `task-id` use the scope (as already said in the branches section).
2. In both cases, the type and scope must refer to the branch and the message need to describe the changes realized. This way allows a message to be easier to read.

#### Short commit messages
For commit messages, follow the following naming scheme: `type(TASK-ID): message`.
Here, the type and SCOPE must refer to the branch being used for development. The message should briefly describe the content of that change.

#### Long commit messages
For commit messages that need a more detailed description of what was implemented, we recommend the following pattern:
```
type(TASK-ID): subject

- detailed message statement 01
- detailed message statement 02
- detailed message statement 03
```

#### Practical example
Let's imagine a very simple example, using a short commit message with the `task-id` equals to `RD-000`. And it's a code fix, so:
```
fix(RD-000)
```

For the **subject** is a single line that contains a succinct description of the change, in `RD-000` the change is changing the button color from blue to red. The subject will be like:
```
fix(RD-000): changed the button colors
```

Advices:
* no dot(.) in the end;
* don’t capitalize the first letter;
* use the SIMPLE PAST, like: changed, fixed, corrected, removed, added. Instead of using present continuos;
* always place the proper coding quotes around "code", like: I'm changing the function `setAnyThing` (proper quote is "`").

In the message, if you need to add more details, you can add a commit description:
```
fix(RD-000): changed the button colors

- created a new css file for the buttons
- changed the color of the cancel button from blue to red
```

#### Message footer
All breaking changes have to be mentioned as a breaking change block in the footer, starting with the word "BREAKING CHANGE:" with a space or two newlines.
The rest of the commit message is then the description of the change, justification, and migration notes. Please see the [Breaking Changes](#breaking-changes) section for more information.
```
**BREAKING CHANGE:** isolate scope bindings definition has changed and the inject option for the directive controller injection was removed.
```

#### Referencing issues (tickets)
Closed bugs should be listed on a separate line in the footer prefixed with the "Closes" keyword like this:
```
Closes RD-567
```

or in case of multiple issues:
```
Closes RD-567
Closes RD-568
Closes RD-569
```

#### Complete example
```
fix(RD-567+RD-568+RD-569): changed the button colors

- created a new CSS file for the buttons
- changed the color of the cancel button from blue to red

**BREAKING CHANGE:** isolate scope bindings definition has changed and the inject option for the directive controller injection was removed

Closes RD-567
Closes RD-568
Closes RD-569
```

#### Automatic Changelog?
How to generate a change log with this established standard? You can try this:
```
git log <last tag> HEAD --pretty=format:%s
git log <last release> HEAD --grep feature
git bisect skip $(git rev-list --grep irrelevant <good place> HEAD
```

## Versioning
To improve the predictability of releases, it is recommended to use semantic versioning (Semver) which, in short, divides releases into three types MAJOR, MINOR AND PATCH.
Each release should be categorized into one of three types as follows:

### MAJOR
A MAJOR release (1.0.0 -> 2.0.0) must be made when new functionalities are released that need to add native modules to the application, or functionalities that drastically change the functioning/structure/layout of the application.

### MINOR
A MINOR release (1.0.0 -> 1.1.0) must be made when incremental functionalities are released to a version, without the need for new native modules and also without major structural changes in the application.

### PATCH
A PATCH release (1.0.0 -> 1.0.1) must be made when corrections are released for the application in production, in which case native modules cannot be added/removed and, depending on the case, tools such as the Microsoft App Center can be used (aka Code Push) to avoid app store approval and review time.

## Working with Branches
Since a branch <branch_name> is behind `main`, you can follow the following process to update it with `main`:

 1. Check that all changes made are committed and that nothing is left in `stash`;
 2. Push to the remote branch: `git push origin <your-branch>`;
 3. Return to master branch: `git checkout main`;
 4. Update master: `git pull origin main`;
 5. Go back to your branch: `git checkout <your-branch>`;
 6. Rebase master: `git rebase main`;
 7. If conflicts appear, resolve them and then add resolutions (`git add <fixed-file>`), and after resolving all conflicts continue the rebase (`git rebase --continue`);
 8. After the rebase you now need to update your remote branch: `git push -f` (yes `-f`, because after the rebase your commits had their hashes updated);
 9. Now you can finally create the Pull Request (PR) and get some coffee.

## Breaking Changes
When working with libraries, code changes can change the interfaces used by others. These updates must be accompanied by an increment of the library version and the dependency version in the applications that use that library.

In `git`, it's important to flag these **BREAKING CHANGES** in 3 places
1. On commits that might cause breaking changes
    - Example:
    ```
    Commit message: refactor(input): use props for state management
    Commit description: **BREAKING CHANGE:** input behavior now must be implemented by the peer, including value and `handleChange`
    ```
2. In the project's `CHANGELOG.md`
    - Example:
    ```
    - **BREAKING CHANGE:** update input to use props for state management. Input behavior now must be implemented by the peer, including value and `handleChange`
    ```
3. In the pull request of the worked story
    - Example:
    ```
    (...template description/checklist)

    **BREAKING CHANGE:** update input to use props for state management. Input behavior now must be implemented by the peer, including value and `handleChange`
    ```

 ## References

 * [Still using GitFlow? What about a simpler alternative?](https://hackernoon.com/still-using-gitflow-what-about-a-simpler-alternative-74aa9a46b9a3)
 * [Karma Commit Messages](http://karma-runner.github.io/1.0/dev/git-commit-msg.html)
 * [Semantic Versioning](https://semver.org/)
