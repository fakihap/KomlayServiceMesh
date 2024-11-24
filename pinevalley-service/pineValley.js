const express = require('express');
const app = express();
const port = 9090;

app.use(express.json());

const doctorData = {
    "cardiologist": [
        {
            "name": "Dr. John Doe",
            "time": "9:00 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Jane Smith",
            "time": "10:00 AM",
            "hospital": "Pine Valley"
        }
    ],
}

app.post('/pineValley/doctors', (request, response) => {
    const { doctorType } = request.body;
    
    if (!doctorType) {
        return response.status(400).json({
            error: "doctorType is required in request body"
        });
    }

    const doctors = doctorData[doctorType.toLowerCase()] || [];

    response.json({
        doctors: {
            doctor: doctors
        }
    });
});

app.listen(port, () => {
    console.log(`Pine Valley Hospital Service running on port ${port}`);
});