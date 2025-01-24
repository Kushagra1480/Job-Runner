import sqlite3
from datetime import datetime

class DatabaseService:
    def __init__(self, db_path):
        self.conn = sqlite3.connect(db_path)

    def save_entry(self, data):
        cursor = self.conn.cursor()
        company = data[0]
        if '*' in company:
            company.replace('*', '')
        role = data[1]
        url = data[3].split('[link=')[1].split(']')[0]
        current_date = datetime.now()
        date = str(current_date.day) + '/' + str(current_date.month) 
        cursor.execute("""
            INSERT INTO applications 
            (company, position, url, applied_date)
            VALUES (?, ?, ?, ?)
        """, (company, role, url, date))
        self.conn.commit()
        self.conn.close()
        

    def __del__(self):
        if hasattr(self, 'conn'):
            self.conn.close()

