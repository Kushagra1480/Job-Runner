## Job Runner

Job Runner is a command line tool designed to consolidate your software engineering job search in one place. It came from my place to put all my job application data in one place as well as my laziness to pull open job boards. I'm currently pulling jobs from the simplify github repo, I plan on adding more job boards down the line.

### Running Instructions

Optional: Create a virtual environment:
```
python -m venv .
```
Don't forget to activate it before installing the requirements
```
activate 
```

Install the requirements:
```
python install -r requirements.txt
```

Cd into the job_runner folder:
```
cd ./job_runner
```
Run the main script: 
```
python main.py <args>
```
Currently, the arguements are:
- add: add a new application manually 
- list: list all applications
- new: view job listings from the last two days
- export: export application data into a csv file

The Simplify Repo can be found [here](https://github.com/SimplifyJobs/New-Grad-Positions)
