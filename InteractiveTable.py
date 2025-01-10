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

    def render(self):
        os.system('cls' if os.name == 'nt' else 'clear')
        self.console.print(self.create_table())

    def run(self):
        self.render()
        
        while True:
            event = keyboard.read_event()
            
            if event.event_type == 'down':  
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
            time.sleep(0.01)
