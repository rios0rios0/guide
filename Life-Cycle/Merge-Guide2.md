# Merge Guide

### Case 1
![](Life-Cycle/Assets/case-1-commits-before.png)

In this case, where we have a second (third, fourth or N branches), dependent on a first branch, we need to merge the last branch (in this case test 2), before merging the first into main (which would be test 1 ).

The merge order for this case is:

1 - > merge test/2 into test/1
2 -> merge test/1 into main

This could have been repeated infinitely many times, like:

x -> thousands of branches

1 - > merge test/4 into test/3

2 - > merge test/3 into test/2

3 - > merge test/2 into test/1

4 -> merge test/1 into main

Result:

![](Life-Cycle/Assets/case-1-commits-result.png)

--------------------------------
### Case 2
In case we have 2 branches starting from main, the procedure must be done as follows:

1 - merge test 1 into main

2 - rebase test 2 with main

3 - merge test 2 into main

This way, there will be two independent triangles. It's crucial that you don't forget to update "main" locally after the first merge. The procedure can be repeated infinite times. For example:

1 - merge test 1 into main

2 - rebase test 2 with main

3 - merge test 2 into main

4 - rebase test 3 with main

5 - merge test 3 into main

N - ...... thousands of branches

Note: if after updating the main and trying to rebase to the main of the 2nd branch this conflict occurs, go to the files indicated in your IDE and fix the conflicts.

After that, give the commands

-   git rebase –-continue
-   git push –f

![](Life-Cycle/Assets/case-2commits-error.png)

Then update your main and merge. The graph will look like this:
![](Life-Cycle/Assets/case-2-commits-result.png)
