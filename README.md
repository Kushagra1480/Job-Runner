# 🎯 Job Search Tracker CLI 

```
     _       _      _____               _             
    | | ___ | |__  |_   _| __ __ _  ___| | _____ _ __ 
 _  | |/ _ \| '_ \   | || '__/ _` |/ __| |/ / _ \ '__|
| |_| | (_) | |_) |  | || | | (_| | (__|   <  __/ |   
 \___/ \___/|_.__/   |_||_|  \__,_|\___|_|\_\___|_|   
```

[![Python 3.6+](https://img.shields.io/badge/python-3.6+-blue.svg)](https://www.python.org/downloads/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Code style: black](https://img.shields.io/badge/code%20style-black-000000.svg)](https://github.com/psf/black)

> 🚀 A sleek CLI tool to track your job applications and monitor new grad tech positions

## ✨ Features

<details open>
<summary>Core Functionality</summary>

- 📝 **Track Applications**: Manage your job applications with ease
- 🔍 **Monitor New Listings**: Auto-fetch new grad positions from Simplify
- 📊 **Application Status**: Visual status tracking with emojis
- 📤 **Export Data**: Export your applications in CSV/JSON formats
</details>

## 🚀 Quick Start

```bash
# Install dependencies
pip install typer pandas requests beautifulsoup4 rich

# Run the tracker
python job_tracker.py
```

## 🎮 Commands

```
╭─ Commands ────────────────────────────────╮
│ add     Add new application              │
│ list    View your applications           │
│ new     Check new listings               │
│ export  Export your data                 │
╰──────────────────────────────────────────╯
```

## 💫 Usage Examples

```bash
# Add new application
python job_tracker.py add

# List all applications
python job_tracker.py list

# Check new listings from last 3 days
python job_tracker.py new --days 3

# Export data as CSV
python job_tracker.py export --format csv
```

## 🌟 Status Indicators

| Emoji | Status      |
|-------|------------|
| 📝    | Preparing  |
| ✉️    | Applied    |
| 💬    | Interviewing|
| 🎉    | Offered    |
| ❌    | Rejected   |

## 🔧 Data Storage

Applications are stored in a SQLite database at `~/.job_tracker.db`

## 🛠️ Dependencies

- typer: CLI interface
- rich: Terminal formatting
- pandas: Data handling
- requests: API calls
- beautifulsoup4: HTML parsing

---

<p align="center">
Made with ❤️ for job seekers
<br>
Happy job hunting! 🎯
</p>
