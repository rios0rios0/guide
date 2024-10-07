## The BFG
It's a simple way to clean your repository of unwanted data in commits on a historical level. First step is to install it:

## Install the jar file of BFG
```bash
wget https://repo1.maven.org/maven2/com/github/git-tools/bfg/1.14.0/bfg-1.14.0.jar -O bfg.jar
```

## Clone the repository you want to clean
Clone the repository in the working directory
```bash
git clone --mirror git://example.com/<your-repo>.git
```

## Prepare replacements.txt
Now you may also prepare a replacements.txt file, which will contain the secrets/credentials you want to replace across the history of your repository. Here is the form. This secrets you can get by running a security tool of your choice.

## Perform cleaning of the repository
Make sure when running this command, both the repoitory and the replacements.txt file are in the same directory level.
```bash
bfg --replace-text replacements.txt <your-repo>.git
```

## Once the clean has been performed, run the following mandatory commands
```bash
git reflog expire --expire=now --all
```
```bash
git gc --prune=now --aggressive
```

## Now make the push into the repository
```bash
git push --mirror
```
