import keyboard
import os
import time
import webbrowser
from rich.text import Text
from rich.console import Console
from rich.table import Table

class InteractiveTable:
    def __init__(self, data, columns):
        self.data = data
        self.columns = ['#'] + columns
        self.selected_row = 0
        self.console = Console()
        
    def create_table(self):
        table = Table(show_header=True, header_style="bold green")
        for col in self.columns:
            table.add_column(col)
        for idx, row in enumerate(self.data):
            style = "reverse" if idx == self.selected_row else ""
            row_data = [Text(str(idx + 1), style=style)]
            for cell in row:
                if '[link=' in str(cell):
                    row_data.append(cell)
                else:
                    row_data.append(Text(str(cell), style=style))
            table.add_row(*row_data)
        return table

    def print_commands(self):
        self.console.print("\nCommands:", style="bold blue")
        self.console.print("↑/k/nk: Move up n times (or type 'up')")
        self.console.print("↓/j/nj: Move down n times (or type 'down')")
        self.console.print("enter/o/no: Open link at row n")
        self.console.print("oa: Open all links")
        self.console.print("q: Quit")
        self.console.print("\nEnter command: ", end="")

    def render(self):
        os.system('cls' if os.name == 'nt' else 'clear')
        self.console.print(self.create_table())
        self.print_commands()

    def open_link(self, row_idx):
        if 0 <= row_idx <= len(self.data):
            apply_link = str(self.data[row_idx][3])
            if '[link=' in apply_link:
                url = apply_link.split('[link=')[1].split(']')[0]
                webbrowser.open(url)
                return True
        return False
            
    def open_all_links(self):
        opened_count = 0
        for row_idx in range(len(self.data)):
            if self.open_link(row_idx):
                opened_count += 1
        return opened_count

    def handle_command(self, command_string):
        if not command_string:
            return True
        num = ''
        i = 0
        while i < len(command_string) and command_string[i].isdigit():
            num += command_string[i]
            i += 1
        cmd = command_string[i:].lower() if i < len(command_string) else ''
        count = int(num) if num else 1
        if cmd == 'q':
            return False
        elif cmd == 'j':
            self.selected_row = min(self.selected_row + count, len(self.data) - 1)
        elif cmd == 'k':
            self.selected_row = max(self.selected_row - count, 0)
        elif cmd == 'o':
            if num:
                row_idx = int(num) - 1
                if self.open_link(row_idx):
                    input("\nPress Enter when you're ready to continue...")
            else:
                if self.open_link(self.selected_row):
                    input("\nPress Enter when you're ready to continue...")
        elif cmd == 'oa':
            opened = self.open_all_links()
            self.console.print(f"\nOpened {opened} links.")
            input("\nPress Enter when you're ready to continue...")
        return True

        
    def run(self):
        running = True
        while running:
            self.render()
            command = input()
            running = self.handle_command(command)
