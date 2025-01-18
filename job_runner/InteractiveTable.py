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
        self.columns = columns
        self.selected_row = 0
        self.console = Console()
        
    def create_table(self):
        table = Table(show_header=True, header_style="bold green")
        for col in self.columns:
            table.add_column(col)
        for idx, row in enumerate(self.data):
            style = "reverse" if idx == self.selected_row else ""
            row_data = []
            for cell in row:
                if '[link=' in str(cell):
                    row_data.append(cell)
                else:
                    row_data.append(Text(str(cell), style=style))
            table.add_row(*row_data)
        return table

    def print_commands(self):
        self.console.print("\nCommands:", style="bold blue")
        self.console.print("↑/k: Move up (or type 'up')")
        self.console.print("↓/j: Move down (or type 'down')")
        self.console.print("enter/o: Open link (or type 'open')")
        self.console.print("q: Quit")
        self.console.print("\nEnter command: ", end="")

    def render(self):
        os.system('cls' if os.name == 'nt' else 'clear')
        self.console.print(self.create_table())
        self.print_commands()

    def handle_command(self, command_string):
        i = 0
        while i < len(command_string):
            current_cmd = command_string[i].lower()

            if current_cmd == 'q':
                return False
            elif current_cmd == 'j':
                moves = 1
                while i + 1 < len(command_string) and command_string[i + 1].lower() == 'j':
                    moves += 1
                    i += 1
                self.selected_row = min(self.selected_row + moves, len(self.data) - 1)
            elif current_cmd == 'k':
                moves = 1
                while i + 1 < len(command_string) and command_string[i + 1].lower() == 'k':
                    moves += 1
                    i += 1
                self.selected_row = max(self.selected_row - moves, 0)
            elif current_cmd == 'o':
                apply_link = str(self.data[self.selected_row][3])
                if '[link=' in apply_link:
                    url = apply_link.split('[link=')[1].split(']')[0]
                    webbrowser.open(url)
                    input("\nPress Enter when you're ready to continue...")
            i += 1
        return True

            

    def run(self):
        running = True
        while running:
            self.render()
            command = input()
            running = self.handle_command(command)
