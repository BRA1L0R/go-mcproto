# **Contributing to go-mcproto**

You can contribute by:

- Reporting a bug<br>
- Discussing the current state of the code<br>
- Submitting a fix<br>
- Proposing new features<br>
- Becoming a maintainer<br>

## **How to setup the repository**

1. Run <code>git clone https://github.com/BRA1L0R/go-mcproto</code>
2. Install golangcgi-lint, as it is required and run by git hooks. Follow [this guide](https://golangci-lint.run/usage/install/#local-installation)
   2.1. I also recommend setting up editor integration (most of them are supported) by following [this guide](https://golangci-lint.run/usage/install/#local-installation)
3. Run <code>./init-hooks.sh</code> to setup git hooks for unit testing and linting

## **Contributing guidelines**

1. If any features are added, they need to be unit tested.<br>
2. The pull request must contain a comprehensive description where you explain why you are adding a certain feature or why you are modifying certain behaviours of the code.<br>
3. The code must have no warnings caused by golangcgi-lint and no errors in general.<br>
4. Don't create pull requests that have low effort, like modifying the README.md.<br>

## **Bug reports**

When you report a bug it is important that you follow these simple rules! This will help me fix the bug faster!<br>

- Write a quick summary of what happend.
- Put your code and your output in the description of the bug.
- Write the excpected behavior and what actually happens.
- Write down the machine configuration that you are experiencing this bug on, including O.S, CPU, RAM, SSD/HDD, GPU etc.

## New to github? No problem, follow this little guide!

1. Fork the official repository

2. Clone YOUR fork with:<br>

<code>git clone https://github.com/yourusername/repository.git</code><br>

3. Navigate to your local repository with:<br>

<code>cd go-mcproto</code><br>

4. Make changes in your local repository by creating, editing files etc.

5. Cofirm your changes with:<br>

<code>git add -A</code> to confirm every change made<br>

<code>git add file1 file2</code> to confirm only the specified files<br>

6. Commit them with:<br>

<code>git commit -m "Description of the changes"</code><br>

6. Push your changes to your fork with<br>

<code>git push</code><br>

7. Begin the pull request by going to your forked repository and clicking on "Contribute"<br>

![image](https://i.ibb.co/QFV8Wwb/image.png)

8. Enter the edit form by clicking the green button<br>

![image](https://i.ibb.co/1bDbN7P/image.png)

9. Write the title and a comment and finally click the green button to create the pull request
   ![image](https://i.ibb.co/wCQppy7/image.png)
