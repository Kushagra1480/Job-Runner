#!/usr/bin/env python3

import typer
import sqlite3
from typing import Optional
from datetime import datetime
import requests
import base64
import webbrowser
from rich.console import Console
from rich.table import Table
from rich.progress import Progress, SpinnerColumn, TextColumn
from rich.panel import Panel
from rich.prompt import Prompt
from rich.live import Live
import pandas as pd
from pathlib import Path
from bs4 import BeautifulSoup
from rich.text import Text
import keyboard
import os
import time


class InteractiveTable:
    def __init__(self, data, columns):
        self.data = data
        self.columns = columns
        self.selected_row = 0
        self.console = Console()
        
    def create_table(self):
        table = Table(show_header=True, header_style="bold green")
        
        # Add columns
        for col in self.columns:
            table.add_column(col)
        
        # Add rows with highlighting for selected row
        for idx, row in enumerate(self.data):
            style = "reverse" if idx == self.selected_row else ""
            row_data = []
            for cell in row:
                if '[link=' in str(cell):
                    # Keep link formatting as is
                    row_data.append(cell)
                else:
                    row_data.append(Text(str(cell), style=style))
            table.add_row(*row_data)
        
        return table

    def render(self):
        """Clear screen and render table"""
        os.system('cls' if os.name == 'nt' else 'clear')
        self.console.print(self.create_table())

    def run(self):
        self.render()
        
        while True:
            event = keyboard.read_event()
            
            if event.event_type == 'down':  # Only handle key press events
                if event.name == 'q':
                    break
                elif event.name == 'up' and self.selected_row > 0:
                    self.selected_row -= 1
                    self.render()
                elif event.name == 'down' and self.selected_row < len(self.data) - 1:
                    self.selected_row += 1
                    self.render()
                elif event.name == 'enter':
                    apply_link = str(self.data[self.selected_row][3])
                    if '[link=' in apply_link:
                        url = apply_link.split('[link=')[1].split(']')[0]
                        webbrowser.open(url)
                        self.render()
            
            # Small delay to prevent high CPU usage
            time.sleep(0.01)
app = typer.Typer(help="Track your job applications and monitor new listings")
console = Console()

# Database setup
DB_PATH = Path.home() / ".job_tracker.db"

def init_db():
    """Initialize the SQLite database with required tables"""
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    # Applications table
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS applications (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        company TEXT NOT NULL,
        position TEXT NOT NULL,
        source TEXT DEFAULT 'manual',
        status TEXT DEFAULT 'applied',
        applied_date DATETIME NOT NULL,
        last_updated DATETIME NOT NULL,
        url TEXT,
        notes TEXT,
        favorite BOOLEAN DEFAULT 0
    )
    """)
    
    # Last fetch tracking for job boards
    cursor.execute("""
    CREATE TABLE IF NOT EXISTS last_fetch (
        source TEXT PRIMARY KEY,
        last_timestamp DATETIME NOT NULL
    )
    """)
    
    conn.commit()
    conn.close()

def show_welcome():
    """Display welcome message and ASCII art banner"""
    welcome_text = """
    🎯 [bold blue]Job Search Tracker[/bold blue] 🎯
    Track applications & monitor new listings
    
    Commands:
    • add - Add new application
    • list - View your applications
    • new - Check new listings
    • export - Export your data
    """
    console.print(Panel.fit(welcome_text, border_style="bright_blue"))

@app.command()
def add(
    company: str = typer.Option(..., prompt=True),
    position: str = typer.Option(..., prompt=True),
    url: Optional[str] = typer.Option(None, prompt="Job posting URL (optional)"),
    notes: Optional[str] = typer.Option(None, prompt="Any notes (optional)")
):
    """Add a new job application with interactive prompts"""
    status = Prompt.ask(
        "Application status",
        choices=["preparing", "applied", "interviewing", "offered", "rejected"],
        default="applied"
    )
    
    with Progress(
        SpinnerColumn(),
        TextColumn("[progress.description]{task.description}"),
        transient=True,
    ) as progress:
        task = progress.add_task("Adding application...", total=None)
        
        conn = sqlite3.connect(DB_PATH)
        cursor = conn.cursor()
        now = datetime.now()
        
        cursor.execute("""
            INSERT INTO applications 
            (company, position, status, url, notes, applied_date, last_updated)
            VALUES (?, ?, ?, ?, ?, ?, ?)
        """, (company, position, status, url, notes, now, now))
        
        conn.commit()
        conn.close()
    
    console.print(f"✅ Added application for [bold]{position}[/bold] at [bold]{company}[/bold]")

@app.command()
def list(
    status: Optional[str] = typer.Option(None, help="Filter by status"),
    favorite: bool = typer.Option(False, help="Show only favorited applications")
):
    """List all tracked applications in a formatted table"""
    conn = sqlite3.connect(DB_PATH)
    cursor = conn.cursor()
    
    query = "SELECT * FROM applications"
    params = []
    
    if status:
        query += " WHERE status = ?"
        params.append(status)
    if favorite:
        query += " WHERE favorite = 1"
    
    query += " ORDER BY applied_date DESC"
    
    cursor.execute(query, params)
    applications = cursor.fetchall()
    
    if not applications:
        console.print("No applications found! 📝")
        return
    
    table = Table(show_header=True, header_style="bold magenta")
    table.add_column("Company")
    table.add_column("Position")
    table.add_column("Status", style="cyan")
    table.add_column("Applied Date")
    table.add_column("Notes")
    
    status_emoji = {
        "preparing": "📝",
        "applied": "✉️",
        "interviewing": "💬",
        "offered": "🎉",
        "rejected": "❌"
    }
    
    for app in applications:
        status = f"{status_emoji.get(app[3], '•')} {app[3]}"
        applied_date = datetime.fromisoformat(app[5]).strftime("%Y-%m-%d")
        notes = (app[8][:30] + "...") if app[8] and len(app[8]) > 30 else (app[8] or "")
        
        table.add_row(
            app[1],  # company
            app[2],  # position
            status,
            applied_date,
            notes
        )
    
    console.print(table)
    
    # Print summary statistics
    cursor.execute("SELECT status, COUNT(*) FROM applications GROUP BY status")
    stats = cursor.fetchall()
    
    console.print("\n[bold]Summary:[/bold]")
    for status, count in stats:
        console.print(f"{status_emoji.get(status, '•')} {status}: {count}")
    
    conn.close()

@app.command()
def new(days: int = typer.Option(2, help="Show listings from last N days")):
    """Check for new listings from Simplify's GitHub repo"""
    with Progress() as progress:
        fetch_task = progress.add_task("Fetching new listings...", total=100)
        
        # GitHub API endpoint for Simplify's README
        url = "https://api.github.com/repos/SimplifyJobs/New-Grad-Positions/contents/README.md?ref=dev"
        
        progress.update(fetch_task, advance=30)
        response = requests.get(url)
        content = base64.b64decode(response.json()["content"]).decode()
        
        progress.update(fetch_task, advance=30)
        
        # Find the table in the content
        lines = content.split('\n')
        
        # Look for the listings table
        table_start = None
        for i, line in enumerate(lines):
            if '| Company' in line:
                table_start = i
                break
        
        if table_start is None:
            console.print("❌ Couldn't find listings table")
            return
            
        # Get table lines
        table_lines = []
        for line in lines[table_start:]:
            if not line.startswith('|'):
                break
            table_lines.append(line)
        
        if not table_lines:
            console.print("❌ No listings found in the expected format")
            return
        
        # Convert to DataFrame for easier handling
        df = pd.DataFrame([
            [cell.strip() for cell in line.split('|')[1:-1]]
            for line in table_lines[2:]  # Skip header and separator
        ])


        # Add column names
        df.columns = ['company', 'position', 'location', 'url', 'date_posted']
        
        

        # Convert dates and filter for recent posts
        from datetime import datetime, timedelta
        
        def parse_date(date_str):
            try:
                date = datetime.strptime(date_str, "%b %d")
                current_date = datetime.now()

                year = current_date.year - 1 if date.month > current_date.month else current_date.year
                
                return datetime(year, date.month, date.day)

            except:
                return None
        
        cutoff_date = datetime.now() - timedelta(days=days)

        
        df['date_posted'] = df['date_posted'].apply(parse_date)
        
        df = df[df['date_posted'] >= cutoff_date].copy()
        
        df = df.sort_values('date_posted', ascending=False)
        
        progress.update(fetch_task, advance=40)
    
    if df.empty:
        console.print("No new listings found! 🔍")
        return
    def extract_first_href(html):
        try:
            soup = BeautifulSoup(html, 'html.parser')
            link = soup.find('a')
            return link['href'] if link else None
        except Exception as e:
            return None
        
    def get_company_name(markdown):
        try:
            text = markdown.strip('*')
            start = text.find('[') + 1
            end = text.find(']')
            if start > 0 and end > start:
                return text[start:end]
            return markdown
        except:
            return markdown

    df['company'] = df['company'].apply(get_company_name)
    df['url'] = df['url'].apply(extract_first_href)
    table_data = []
    

    # Display new listings
    table = Table(show_header=True, header_style="bold green")
    table.add_column("Company")
    table.add_column("Position")
    table.add_column("Location")
    table.add_column("Link")
    table.add_column("Posted", style="cyan")


    console.print("\n[bold]Latest Listings:[/bold]")
    for _, row in df.iterrows():
        if not pd.isna(row['url']):
            apply_link = f"[link={row['url']}][#0000FF]🔗 Apply[/#0000FF][/link]"
        else:
            apply_link = "🔒" 
        table_data.append([
                row['company'],
                row['position'],
                row['location'],
                apply_link,
                row['date_posted'].strftime("%Y-%m-%d")
            ]) 
    interactive_table = InteractiveTable(table_data, df.columns)
    interactive_table.run()


    """ with Live(table, refresh_per_second=4) as live:
        for _, row in df.iterrows():
            if not pd.isna(row['url']):
                apply_link = f"[link={row['url']}][#0000FF]🔗 Apply[/#0000FF][/link]"
            else:
                apply_link = "🔒" 

            table.add_row(
                row['company'],
                row['position'],
                row['location'],
                apply_link,
                row['date_posted'].strftime("%Y-%m-%d")
            ) """
    
    

@app.command()
def export(
    format: str = typer.Option("csv", help="Export format (csv/json)"),
    output: Optional[str] = typer.Option(None, help="Output file path")
):
    """Export your application data"""
    conn = sqlite3.connect(DB_PATH)
    df = pd.read_sql_query("SELECT * FROM applications", conn)
    conn.close()
    
    if not output:
        output = f"job_applications_{datetime.now().strftime('%Y%m%d')}.{format}"
    
    with Progress(
        SpinnerColumn(),
        TextColumn("[progress.description]{task.description}"),
        transient=True,
    ) as progress:
        progress.add_task("Exporting data...", total=None)
        
        if format == "csv":
            df.to_csv(output, index=False)
        else:
            df.to_json(output, orient="records", indent=2)
    
    console.print(f"✅ Exported data to [bold]{output}[/bold]")

@app.callback(invoke_without_command=True)
def main(
    ctx: typer.Context,
    help: bool = typer.Option(False, "--help", "-h", is_eager=True)
):
    """Initialize app and show welcome message"""
    if not DB_PATH.exists():
        init_db()
    
    if ctx.invoked_subcommand is None:
        show_welcome()
            
if __name__ == "__main__":
    app()