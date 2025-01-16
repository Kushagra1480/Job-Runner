from setuptools import setup, find_packages

setup(
    name="job_runner",
    version="0.1",
    packages=find_packages(),
    install_requires=[
        'typer>=0.9.0',
        'rich>=13.7.0',
        'pandas>=2.2.0',
        'requests>=2.31.0',
    ],
    entry_points={
        'console_scripts': [
            'jrun=job_runner.main:app',
        ],
    },
)