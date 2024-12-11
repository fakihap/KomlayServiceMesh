from flask import Flask
from flask_restful import Resource, Api
import json
from faker import Faker

app = Flask(__name__)
api = Api(app)
fake = Faker()

nData = 5
doctorList = {}
doctorType = ["surgeon", "dentist", "neurosurgeon"]

for type in doctorType:
    doctors = []
    for _ in range(nData):
        data = {}
        data['name'] = fake.name()
        data['time'] = fake.time()
        doctors.append(data)

    doctorList[type] = {
        "our_doctors": doctors
    }

class Doctors(Resource):
    def get(self, doctorType):
        return doctorList[doctorType]
    
api.add_resource(Doctors, '/grandOak/doctors/<doctorType>')

if __name__ == '__main__':
    app.run(debug=True)