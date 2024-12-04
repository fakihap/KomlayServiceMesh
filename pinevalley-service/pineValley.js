const express = require('express');
const app = express();
const port = 9090;

app.use(express.json());

const doctorData = {
    "surgeon": [
        {
            "name": "Dr. Ava Sinclair",
            "time": "3:15 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Leo Montgomery",
            "time": "10:00 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Mia Carter",
            "time": "04:15 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Ethan Walker",
            "time": "10:09 PM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Lily Bennett",
            "time": "6:58 AM",
            "hospital": "Pine Valley"
        }
    ],
    "dentist": [
        {
            "name": "Dr. Jack Harper",
            "time": "11:34 PM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Grace Turner",
            "time": "04:17 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Noah Parker",
            "time": "08:15 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Chloe Matthews",
            "time": "02:19 PM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Mason Reed",
            "time": "09:52 AM",
            "hospital": "Pine Valley"
        }
    ],
    "neurosurgeon": [
        {
            "name": "Dr. Zoe Phillips",
            "time": "01:25 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Ryan Foster",
            "time": "11:02 PM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Ella Mitchell",
            "time": "11:25 AM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Lucas Gray",
            "time": "02:14 PM",
            "hospital": "Pine Valley"
        },
        {
            "name": "Dr. Sophia Collins",
            "time": "01:12 PM",
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