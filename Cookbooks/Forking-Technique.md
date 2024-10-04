## Context
When we need to fork some project, we should follow some standards to be up-to-date with the world outside.
Assuming that a fork is primarily made from an open-source project, it’s common that the maintainers are often too busy and take too long to review and merge the PR's from the community.
Even for tiny changes such as adding an icon on the theme.

To solve this kind of situation, we’ve done this standard to guide this.
We intend to keep the compatibility of these and re-base all the changes of newer community versions.

## Technique
To keep it all organized, we’ve defined some patterns:
* The `main` branch is reserved for the community version of the projects. Inside the forks, we are always working in the `custom` branch, which works as a `main` for us.
* To keep our forks up to date with the base project and also apply our changes, we’re using a strategy of always re-base the `main` (community) branch, and placing our customization on top of the newest version.
* Our version tags are going to be an incremental number over the community version.
For example, if the community is on version `1.0.0`, and we are synced with them, then our fork version would be `1.0.0.0`.
If we release a new version of our fork, we increase to `1.0.0.1`. This last digit is incremental and it resets with the community version, so, if the community goes for `1.0.1`, and we re-base the new version, we are on `1.0.1.1`.
