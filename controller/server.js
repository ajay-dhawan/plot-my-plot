// Import required modules
import express from 'express';
import mongoose from 'mongoose';
import { healthCheckPlus } from 'health-check-plus';
import dotenv from 'dotenv';

// Initialize Express app
const app = express();
const port = 3010;

// Load environment variables from .env file
dotenv.config();

// Middleware to parse JSON bodies
app.use(express.json());

// Middleware for health checks
app.use(healthCheckPlus);

// Connect to MongoDB
mongoose.connect(process.env.MONGO_DB, {})
  .then(() => {
    console.log('Connected to MongoDB');
  })
  .catch((err) => {
    console.error('Error connecting to MongoDB:', err);
  });

// Define mongoose schema for data
const dataSchema = new mongoose.Schema({
  CPUUtilizationPercentage: Number,
  CPUError: String,
  FreeMemory: Number,
  TotalMemory: Number,
  MemoryError: String,
  TimeStamp: String,
  db_ts: { type: Date, default: Date.now }
}, {
  versionKey: false
});

// Create mongoose model based on schema
const Data = mongoose.model('Data', dataSchema);

// POST endpoint to save system data
app.post('/system-data', async (req, res) => {
  const newData = new Data({
    CPUUtilizationPercentage: req.body.CPUUtilizationPercentage,
    CPUError: req.body.CPUError,
    FreeMemory: req.body.FreeMemory,
    TotalMemory: req.body.TotalMemory,
    MemoryError: req.body.MemoryError,
    TimeStamp: req.body.TimeStamp
  });

  console.log(newData);

  try {
    // Save data to MongoDB
    await newData.save();
    res.status(201).send('Data saved successfully');
  } catch (err) {
    // Handle error if data save fails
    res.status(400).send(err);
  }
});

// Start the server
app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
