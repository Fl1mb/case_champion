import service
import database
import time

time.sleep(3)

database.create_tables()
service.serve()