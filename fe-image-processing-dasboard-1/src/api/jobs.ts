import axios from 'axios';
import { Job } from '../types';

// Base API URL
const API_URL = 'http://localhost:8080';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * Fetch all image processing jobs
 */
export const fetchJobs = async (): Promise<Job[]> => {
  try {
    const response = await api.get('/jobs');
    return response.data.data;
  } catch (error) {
    console.error('Error fetching jobs:', error);
    throw error;
  }
};

/**
 * Fetch jobs filtered by status
 */
export const fetchJobsByStatus = async (status: string): Promise<Job[]> => {
  try {
    const response = await api.get(`/jobs/status/${status}`);
    return response.data;
  } catch (error) {
    console.error('Error fetching jobs by status:', error);
    throw error;
  }
};

/**
 * Upload multiple images and create processing jobs
 */
export const uploadImages = async (files: File[]): Promise<void> => {
  try {
    const formData = new FormData();
    files.forEach(file => {
      formData.append('images', file);
    });
    
    await api.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  } catch (error) {
    console.error('Error uploading images:', error);
    throw error;
  }
};

/**
 * Retry a failed job
 */
export const retryJob = async (id: number): Promise<void> => {
  try {
    await api.post(`/jobs/${id}/retry`);
  } catch (error) {
    console.error(`Error retrying job ${id}:`, error);
    throw error;
  }
};