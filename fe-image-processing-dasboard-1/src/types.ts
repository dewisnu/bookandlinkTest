export type JobStatus = 'pending' | 'processing' | 'completed' | 'failed';

export interface Job {
  id: number;
  filename: string;
  original_size?: number;
  compressed_size?: number;
  compressed_file_name?: string;
  status: JobStatus;
  error_message?: string;
  created_at?: string;
  updated_at?: string;
}