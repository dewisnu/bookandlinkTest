-- Create database schema for image processing pipeline

-- Create the image_jobs table
CREATE TABLE IF NOT EXISTS image_jobs (
  id SERIAL PRIMARY KEY,
  filename VARCHAR(255) NOT NULL,
  original_size BIGINT,
  compressed_size BIGINT,
  compressed_url TEXT,
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  error_message TEXT,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Add index for status for faster queries when filtering by status
CREATE INDEX IF NOT EXISTS idx_image_jobs_status ON image_jobs (status);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger to automatically update updated_at on update
DROP TRIGGER IF EXISTS trigger_update_image_jobs_updated_at ON image_jobs;
CREATE TRIGGER trigger_update_image_jobs_updated_at
BEFORE UPDATE ON image_jobs
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();