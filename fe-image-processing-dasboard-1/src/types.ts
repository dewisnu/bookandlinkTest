export type JobStatus = 'pending' | 'processing' | 'complete' | 'failed';

export interface Job {
  id: number;
  filename: string;
  original_size?: number;
  compressed_size?: number;
  compressed_url?: string;
  status: JobStatus;
  error_message?: string;
  created_at?: string;
  updated_at?: string;
}