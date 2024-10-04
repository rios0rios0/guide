## Context
This standard is intended to guide how to proceed in creating files, naming functions, and variables.
If there's no override documentation (languages standards inside the same level of this document), this standard must be followed as primary instructions.
The layers mentioned below are defined and exemplified inside this page [here](Life-Cycle/Architecture/Backend-Design.md).

## File Structure
The default file structure is described in the [backend design section](Life-Cycle/Architecture/Backend-Design.md).

## Services
If you need to create a listener for the service, we usually use the same names as the listeners for the commands.
When you are naming this file you MUST be careful considering the entities and if they are operated in batch.

1. **Note:** the operation is the kind of "thing" you're doing. Like: `List`, `Get`, `Delete`, `Insert`, `Update`, `BatchDelete`, `BatchInsert`, `BatchUpdate`, `DeleteAll` and so on.
